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
		token, refreshToken, _ := generate.TokenGenerator(*foundUser.Email,*foundUser.Nama_Depan, *foundUser.User_ID)
		generator.UpdateAllTokens(token, refreshToken, foundUser.User_ID)
		c.JSON(http.StatusFound, foundUser)

	}
}
func ProductViewAdmin() gin.HandlerFunc {

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
		if data.Err() == nil{
			log.Println(err)
			c.JSON(400, "ga valid")
			return
		}

		defer cancel()
		c.JSON(200, listProduk)
		
	}
}

func SearchProductByQuery() gin.HandlerFunc{

}