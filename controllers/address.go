package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilhamadikusuma31/golang-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TambahAlamat() gin.HandlerFunc {
	return func(c *gin.Context) {
        userQueryID := c.Query("id")
		if userQueryID == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			c.Abort()
			return
		
		userID, err := primitive.ObjectIDFromHex(queryUserID)
		if err!= nil {
			c.JSON(500,"internal server error")
		}

		var alamats models.Alamat
		alamats.Alamat_ID = primitive.NewObjectID()
		if err = c.BindJSON(&alamats); err != nil {
			c.JSON(http.StatusNotAcceptable, err.Error())
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		filter := bson.D{{ Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userID}} }}
		unwind := bson.D{{ Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value:"$alamat"}} }}
		group  	:= bson.D{{ Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$addres_id"}} }}
}}
func EditAlamatRumah() gin.HandlerFunc {

}
func EditAlamatKantor() gin.HandlerFunc {

}
func HapusAlamat() gin.HandlerFunc {
	return func(c *gin.Context)  {
		queryUserID := c.Query("id")

		if queryUserID  == " "{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error":"invalid penncarian indeks"})
			c.Abort()
			return
		}

		alamats := make([]models.Alamat,0)
		userID, err := primitive.ObjectIDFromHex(queryUserID )

		if err != nil{
			c.JSON(500, "internal server error")
		}

		//buat timeout dengan context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		
		//cari data yang sesuai id yang di req user
		//d => unordered
		filter := bson.D{primitive.E{Key: "_id",Value: userID}}

		//buat update-an datanya ditimpa dengan alamat kosong
		//$set => ini untuk update dari mongo
		update := bson.D{{ Key: "$set", Value: bson.D{primitive.E{Key: "alamat", Value: alamats}} }}

		_,err = UserCollection.UpdateOne(ctx, filter, update)
		if err!=nil{
			c.JSON(404, "salah perintah")
			return
		}

		defer cancel()
		ctx.Done()

		c.JSON(200, "sukses menghapus data")

	}
}