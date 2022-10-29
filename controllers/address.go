package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilhamadikusuma31/golang-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TambahAlamat() gin.HandlerFunc {
	return func(c *gin.Context) {
        userQueryID := c.Query("id")
		if userQueryID == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			c.Abort()
			return
		}
		
		//bikin data baru tapi baru id nya aja
		alamat, err := primitive.ObjectIDFromHex(userQueryID)
		if err!= nil { 
			c.JSON(500,"internal server error")
		}

		var alamats models.Alamat

		//bikin id baru 
		alamats.Alamat_ID = primitive.NewObjectID()

		//semua data alamat yang ada di user tertentu saat ini di bind ke var alamats
		if err = c.BindJSON(&alamats); err != nil {
			c.JSON(http.StatusNotAcceptable, err.Error())
		}

		//bikin timeout dengan context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)


		//bikin aggregate 
		filter := bson.D{{ Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value:alamat}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$alamat"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$Alamat_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}


		//kurson
		kursor, err :=UserCollection.Aggregate(ctx,mongo.Pipeline{filter,unwind, group})
		if err!= nil {
			c.JSON(500, "internal server error")
		}

		//alamat dari suatu user
		var infoAlamat []bson.M
		if err = kursor.All(ctx, &infoAlamat); err != nil {
			panic(err)
		}

		//mengecek jumlah alamat yang dimiliki user
		var jumlahAlamat int
		for _, j := range infoAlamat {
			count := j["total"]
			jumlahAlamat = count.(int)
		}

		//jika jumlah alamat yang dimiliki user tidak melebihi 2 masih bisa ditambah alamat baru
		if jumlahAlamat < 2{
			filter := bson.D{primitive.E{Key:"_id", Value: alamat}} 
			update := bson.D{{ Key: "$push", Value: bson.D{primitive.E{Key: "alamat", Value: alamats}}}}
			_,err := UserCollection.UpdateOne(ctx, filter, update)

			if err != nil {
				fmt.Println(err)
			}
		}else{
			c.JSON(400,"ga bisa")		
		}

		defer cancel()
		ctx.Done()


}}




func EditAlamatRumah() gin.HandlerFunc {
	return func(c *gin.Context)  {
		queryUserID := c.Query("id")

		if queryUserID  == " "{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error":"invalid penncarian indeks"})
			c.Abort()
			return
		}
		userID, err := primitive.ObjectIDFromHex(queryUserID )

		if err != nil{
			c.JSON(500, "internal server error")
		}

		//buat timeout dengan context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//cari data yang sesuai id yang di req user
		//D => ordered
		var editalamats models.Alamat 
		if err:= c.BindJSON(&editalamats); err != nil{
			c.IndentedJSON(400, "gabisa bind alamat")
		}
		
		filter := bson.D{primitive.E{Key: "_id",Value: userID}}
		update := bson.D{{ Key:"$set", Value: bson.D{primitive.E{Key: "alamat.0.nama_rumah", Value: editalamats.Rumah}, {Key:"alamat.0.nama_jalan",Value: editalamats.Jalan}, {Key:"alamat.0.nama_kota",Value: editalamats.Kota}, {Key:"alamat.0.nama_kodepos",Value: editalamats.Kodepos}}}}

		_,err = UserCollection.UpdateOne(ctx, filter, update)
		if err!= nil {
			c.IndentedJSON(500, "ga bisa update alamat rumah, ada yang salah")
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "berhasil mengupdate rumah")
	}

}
func EditAlamatKantor() gin.HandlerFunc {
	return func(c *gin.Context)  {
		queryUserID := c.Query("id")

		if queryUserID  == " "{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error":"invalid penncarian indeks"})
			c.Abort()
			return
		}
		userID, err := primitive.ObjectIDFromHex(queryUserID )

		if err != nil{
			c.JSON(500, "internal server error")
		}

		//buat timeout dengan context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//cari data yang sesuai id yang di req user
		//D => ordered
		var editalamats models.Alamat 
		if err:= c.BindJSON(&editalamats); err != nil{
			c.IndentedJSON(400, "gabisa bind alamat")
		}
		
		filter := bson.D{primitive.E{Key: "_id",Value: userID}}
		update := bson.D{{ Key:"$set", Value: bson.D{primitive.E{Key: "alamat.1.nama_rumah", Value: editalamats.Rumah}, {Key:"alamat.1.nama_jalan",Value: editalamats.Jalan}, {Key:"alamat.1.nama_kota",Value: editalamats.Kota}, {Key:"alamat.1.nama_kodepos",Value: editalamats.Kodepos}}}}

		_,err = UserCollection.UpdateOne(ctx, filter, update)	
		if err!= nil {
			c.IndentedJSON(500, "ga bisa update alamat kantor, ada yang salah")
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "berhasil mengupdate kantor")
	}

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