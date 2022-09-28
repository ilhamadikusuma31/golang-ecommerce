package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilhamadikusuma31/golang-ecommerce/database"
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
		
		userID := c.Query("userID")
		if userID == " "{
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

		err = database.TambahProdukKeKeranjang(ctx, a.produkKoleksi, a.penggunaKoleksi, productID, userID)
		if err!=nil{
			c.JSON(http.StatusInternalServerError, err)
		}

		c.JSON(200, "sukses menambahkan ke keranjang ")

	}
}


func HapusItem() gin.HandlerFunc{

}
func DapatkanItemDariKeranjang() gin.HandlerFunc{

}
func BeliDariKeranjang() gin.HandlerFunc{

}
func BeliCepat() gin.HandlerFunc{

}