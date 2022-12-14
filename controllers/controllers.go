package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ilhamadikusuma31/golang-ecommerce/database"
	"github.com/ilhamadikusuma31/golang-ecommerce/models"
	generate "github.com/ilhamadikusuma31/golang-ecommerce/tokens"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client,"Users")
var ProductCollection *mongo.Collection = database.ProdukData(database.Client, "Products")
var Validate = validator.New()

func HashPassword(pw string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw),14)
	if err!=nil{
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPW string, givenPW string) (bool, string) {
	valid := true
	msg   := ""
	err   := bcrypt.CompareHashAndPassword([]byte(givenPW), []byte(userPW))
	if err!=nil{
		msg = "uname atau pw salah"
		valid = false
	}

	return valid, msg
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second) //ngasih timeout
		defer cancel() //berasal dari line sebelumnya, defer: delay fungsi ini sampai ada fungsi terdekat yg dipanggil

		var pengguna models.User

		if err := c.BindJSON(&pengguna);err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//mem-validasi inputan sesuai dengan struct yang udah dibikin
		validationErr := Validate.Struct(pengguna)
		if validationErr!=nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":validationErr.Error()})
			return
		}


		//cek apakah email sudah ada di DB
		count, err := UserCollection.CountDocuments(ctx, bson.M{"email":pengguna.Email})
		if err!=nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":err})
			return
		}
		if count > 0{
			c.JSON(http.StatusBadRequest, gin.H{"error":"udah ada emailnya bos"})	
			return
		}


		//cek apakah HP sudah ada di DB
		count, err = UserCollection.CountDocuments(ctx, bson.M{"hp":pengguna.Hp})
		defer cancel()
		if err!=nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":err})
			return
		}
		if count > 0{
			c.JSON(http.StatusBadRequest, gin.H{"error":"udah ada no.HP bos"})
			return
		}


		//pengisian value untuk pengguna baru dan di simpan ke DB
		password := HashPassword(*pengguna.Password) //hash password
		pengguna.Password = &password
		pengguna.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339)) //ini kalo di laravel terbuat otomatis created_at sama updated_at
		pengguna.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		pengguna.ID = primitive.NewObjectID() //membuat id baru otomatis
		pengguna.User_ID = pengguna.ID.Hex()


		token, refreshtoken, _ := generate.TokenGenerator(*pengguna.Email, *pengguna.Nama_Depan, *pengguna.Nama_Belakang, pengguna.User_ID)
		pengguna.Token = &token
		pengguna.Refresh_Token = &refreshtoken


		pengguna.KeranjangUser = make([]models.ProdukUser, 0)
		pengguna.Alamat_Detail = make([]models.Alamat, 0)
		pengguna.Pesanan_Status = make([]models.Pesanan, 0)



		//insert ke DB
		_, insertErr := UserCollection.InsertOne(ctx, pengguna)
		if insertErr!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"ga bisa buat daftar boskuh" })
			return
		}
		defer cancel()

		c.JSON(http.StatusCreated, "berhasil daftar ke sistem")




	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context)  {
		//membuat timeout
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//manggil class user
		var user models.User
		var foundUser models.User	

		//tempelkan sesuai data request ke class
		err := c.BindJSON(&user)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err})
		} 

		err = UserCollection.FindOne(ctx, bson.M{"email":user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"uname atau password salah"})
			return
		}

		PasswordIsValid, pesan :=VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if !PasswordIsValid{
			c.JSON(http.StatusInternalServerError, gin.H{"error":pesan})
		}

		//membuat token
		token, refreshToken, _ := generate.TokenGenerator(*foundUser.Email,*foundUser.Nama_Depan,*foundUser.Nama_Belakang ,foundUser.User_ID)
		generate.UpdateAllTokens(token, refreshToken, foundUser.User_ID)
		c.JSON(http.StatusFound, foundUser)

	}
}
func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var products models.Produk
		defer cancel()
		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		products.Produk_ID = primitive.NewObjectID()
		_, anyerr := ProductCollection.InsertOne(ctx, products)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Successfully added our Product Admin!!")
	}
}

func SearchProduct() gin.HandlerFunc{
	return func(c *gin.Context)  {
		var listProduk []models.Produk

		//membuat timeout dengan context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()


		//mencari semua data produk
		//bson.d => represntasi terurut berbentuk slice
		data, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil{
			c.JSON(http.StatusInternalServerError, "ada sesuatu yang salah coba lagi nanti boskuh")
			return
		}


		//semua data di pindah valuenya ke var listProduct
		//nb: biar golang tau json
		err = data.All(ctx, &listProduk)
		if err!=nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}


		//close data
		defer data.Close(ctx)
		if data.Err() != nil{
			log.Println(err)
			c.JSON(400, "ga valid")
			return
		}
		defer cancel()

		//feedback data nya untuk dilempar ke FE
		c.JSON(200, listProduk)
		
	}
}

func SearchProductByQuery() gin.HandlerFunc{
	return func(c *gin.Context)  {
		var produkYangDicari []models.Produk
		queryParam := c.Query("name")

		if queryParam == " "{
			log.Println("Query kosong")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error":"invalid indeks yang dicari"})
			c.Abort()
			return
		}

		//bikin timeout dengan context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//cari nama produk sesuai req dari user
		data, err := ProductCollection.Find(ctx, bson.M{"nama_produk":bson.M{"$regex":queryParam}})


		//kalo yang dicar gada yang sesuai
		if err != nil{
			c.JSON(404, "ada kesalaha ketika fetch data")
			return
		}
		defer cancel()


		//data ditempel ke var produkYangDicari
		err = data.All(ctx, &produkYangDicari)
		if err != nil{
			log.Println(err)
			c.JSON(400, "invalid")
			return
		}

		//close data
		defer data.Close(ctx)

		//kalo ada error pas di close
		if data.Err() != nil{
			c.JSON(400, "permintaan invalid")
			return
		}
		defer cancel()

		//feedback data nya untuk dilempar ke FE
		c.JSON(200, produkYangDicari)

	}
}