package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilhamadikusuma31/golang-ecommerce/database"
	"github.com/ilhamadikusuma31/golang-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


type Aplikasi struct{
	produkKoleksi *mongo.Collection
	penggunaKoleksi *mongo.Collection
}

func AplikasiBaru(productColl, userColl *mongo.Collection ) *Aplikasi{
	return &Aplikasi{
		produkKoleksi : productColl,
		penggunaKoleksi : userColl,
	}
}

func (a *Aplikasi) TambahKeKeranjang() gin.HandlerFunc{
	return func (c *gin.Context)  {
		produkQueryId := c.Query("produkID")
		if produkQueryId == " "{
			log.Println("id produk gada")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("id produk gada"))
			return
		}
		
		userQueryID := c.Query("userID")
		if userQueryID == " "{
			log.Println("id user gada")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("id user gada"))
			return
		}

		//bikin id baru
		productID, err := primitive.ObjectIDFromHex(produkQueryId)
		if err!=nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.TambahProdukKeKeranjang(ctx, a.produkKoleksi, a.penggunaKoleksi, productID, userQueryID)
		if err!=nil{
			c.JSON(http.StatusInternalServerError, err)
		}

		c.JSON(200, "sukses menambahkan ke keranjang ")

	}
}


func (a *Aplikasi) HapusItem() gin.HandlerFunc{
	return func (c *gin.Context)  {
		produkQueryId := c.Query("produkID")
		if produkQueryId == " "{
			log.Println("id produk gada")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("id produk gada"))
			return
		}
		
		userQueryID := c.Query("userID")
		if userQueryID == " "{
			log.Println("id user gada")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("id user gada"))
			return
		}

		//bikin id baru
		productID, err := primitive.ObjectIDFromHex(produkQueryId)
		if err!=nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.HapusProdukDariKeranjang(ctx, a.produkKoleksi, a.penggunaKoleksi, productID, userQueryID)

		if err != nil{
			c.JSON(http.StatusInternalServerError,err)
			return
		}

		c.JSON(200,"berhasil menghapus item")

	}
}
func (a *Aplikasi) DapatkanItemDariKeranjang() gin.HandlerFunc{
	return func(c *gin.Context)  {
		
		//tangkap id dari req
		userQueryID := c.Query("id")
	

		//kalo id ga valid
		if userQueryID == " "{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error":"id ga valid"})
			c.Abort()
			return
		}

		//convert type dari hex biar dikenali golang
		userID,_ := primitive.ObjectIDFromHex(userQueryID)
		
		//buat timeout
		ctx, cancel:= context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//cari user yang sesuai id
		var user models.User
		penampung := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userID}})
		err       := penampung.Decode(&user)		
		if err != nil{
			log.Println(err)
			c.JSON(500, "tidak ada user yang ber id tersebut")
			return
		}	


		//aggregate => ada di dokumentasi mongo
		//match => nyari yang sesuai id
		filter := bson.D{{Key: "$match" ,Value: bson.D{primitive.E{Key: "_id", Value: userID}}}}

		//unwind => dipecah per object barangnya dari seseorang user
		unwind := bson.D{{ Key: "$unwind", Value: bson.D{primitive.E{ Key:"path", Value: "keranjang_user" }} }}


		//grouping
		group := bson.D{{ Key: "$group", Value: primitive.E{Key: "_id", Value: "$_id"}}, {Key:"total",Value:bson.D{primitive.E{Key:"$sum",Value:"$keranjang_user.harga"}}} }
		pointer, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter, unwind, group}) 
		if err!=nil{
			log.Println(err)
		}

		//data dilooping dan disimpan di var listing
		var listing []bson.M
		if err = pointer.All(ctx, &listing); err != nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		for _,j := range listing{
			c.JSON(200,j["total"])
			c.JSON(200,user.KeranjangUser)
		}

		ctx.Done()	
	
	}
	
}

func (a *Aplikasi) BeliDariKeranjang() gin.HandlerFunc{
	return func (c *gin.Context)  {
		userQueryID := c.Query("userID")
		if userQueryID == " "{
			log.Println("id user gada")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("id user gada"))
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.BeliProdukDariKeranjang(ctx, a.penggunaKoleksi, userQueryID)
		if err != nil{
			c.JSON(http.StatusInternalServerError, err)
		}

		c.JSON(200,"sukses menyiapkan pesanan")

	}
}
func (a *Aplikasi) BeliCepat() gin.HandlerFunc{
	return func (c *gin.Context)  {
		produkQueryId := c.Query("produkID")
		if produkQueryId == " "{
			log.Println("id produk gada")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("id produk gada"))
			return
		}
		
		userQueryID := c.Query("userID")
		if userQueryID == " "{
			log.Println("id user gada")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("id user gada"))
			return
		}

		//bikin id baru
		productID, err := primitive.ObjectIDFromHex(produkQueryId)
		if err!=nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.PembelianCepat(ctx, a.produkKoleksi, a.penggunaKoleksi, productID, userQueryID)

		if err != nil{
			c.JSON(http.StatusInternalServerError,err)
			return
		}

		c.JSON(200,"berhasil menambahkan pesanan")

	}
}