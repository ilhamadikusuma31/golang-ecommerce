package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func errHandler(e error)  {
	if e!=nil{
		log.Fatal(e)
	}
}

func DBset() *mongo.Client  {
	
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	errHandler(err)

	ctx,cancel := context.WithTimeout(context.Background(), 10*time.Second) //timeout
	defer cancel() 			             //dapet dari line sebelumnya, defer: menjalankan ini tapi delay sampe func terdekat dipanggil

	err = client.Connect(ctx)            //mencoba konek diikutkan dengan timeout
	errHandler(err)
	

	err = client.Ping(context.TODO(),nil) //mencoba nge-ping
	errHandler(err)

	fmt.Println("berhasil konek")
	return client

}

var Client *mongo.Client = DBset()


func UserData(client *mongo.Client, namaKoleksi string) *mongo.Collection {
	var koleksi *mongo.Collection = client.Database("Ecommerce").Collection(namaKoleksi) //db: Ecommerce table:User
	return koleksi
}

func ProdukData(client *mongo.Client, namaKoleksi string) *mongo.Collection {
	var koleksi *mongo.Collection = client.Database("Ecommerce").Collection(namaKoleksi)
	return koleksi
}