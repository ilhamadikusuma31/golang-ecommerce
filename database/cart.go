package database

import (
	"context"
	"errors"
	"log"

	"github.com/ilhamadikusuma31/golang-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindTheProduct = errors.New("ga nemu produknya")
	ErrCantDecodeTheProducts = errors.New("ga nemu produk2nya")
	ErrUserIdIsNotValid = errors.New("user id nya ga valid")
	ErrCantUpdateUser = errors.New("ga bisa update user")
	ErrCantRemoveItemCart = errors.New("ga bisa ngapus item dari keranjang")
	ErrCantGetItem = errors.New("ga bisa dapet item dari keranjang")
	ErrCantBuyCartItem= errors.New("ga bisa beli item yang ada di keranjang")
)


func TambahProdukKeKeranjang(ctx context.Context, prodKoleksi, userKoleksi *mongo.Collection, produkId primitive.ObjectID, userID string) error {
	cariDiDB, err := prodKoleksi.Find(ctx, bson.M{"_id": userID})
	if err!=nil {
		log.Println(err)
		return ErrCantFindTheProduct
	}

	var productCart []models.ProdukUser
	err = cariDiDB.All(ctx, &productCart)
	if err!=nil {
		log.Println(err)
		return ErrCantDecodeTheProducts
	}

	// id produk berbentuk objek id mirip punya mongo
	// sedangkan id user masih berbentuk string
	id, err := primitive.ObjectIDFromHex(userID)
	if err!=nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter:= bson.D{primitive.E{Key: "_id", Value:id}}
	update:= bson.D{primitive.E{Key:"$push", Value:bson.D{primitive.E{Key: "usercart", Value:productCart}}}}

	_, err = userKoleksi.UpdateOne(ctx, filter, update)
	if err!=nil {
		return ErrCantUpdateUser
	} 

	return nil
}

func HapusProdukDariKeranjang(ctx context.Context, prodKoleksi, userKoleksi *mongo.Collection, produkId primitive.ObjectID, userID string) error  {
	id, err := primitive.ObjectIDFromHex(userID)
	if err!=nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter:= bson.D{primitive.E{Key: "_id", Value:id}}
	update:= bson.M{"$pull": bson.M{"usercart":bson.M{"_id":produkId}}}

	_, err = userKoleksi.UpdateMany(ctx, filter, update)
	if err!=nil {
		return ErrCantRemoveItemCart
	} 

	return nil

}
func BeliProdukDariKeranjang()  {
	
}

func PembeliCepat()  {
	
}