package database

import (
	"context"
	"errors"
	"log"
	"time"

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
	update:= bson.D{primitive.E{Key:"$push", Value:bson.D{primitive.E{Key: "keranjang_user", Value:productCart}}}}

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
	update:= bson.M{"$pull": bson.M{"keranjang_user":bson.M{"_id":produkId}}}

	_, err = userKoleksi.UpdateMany(ctx, filter, update)
	if err!=nil {
		return ErrCantRemoveItemCart
	} 

	return nil

}
func BeliProdukDariKeranjang(ctx context.Context, prodKoleksi, userKoleksi *mongo.Collection, produkId primitive.ObjectID, userID string) error   {
	
	id, err := primitive.ObjectIDFromHex(userID)
	if err!=nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var user models.User
	var orderCart models.Pesanan

	orderCart.Pesanan_ID = primitive.NewObjectID()
	orderCart.Keranjang_Pesanan = []models.ProdukUser{}
	orderCart.Orderer_At = time.Now().Local()
	orderCart.Metode_Pembayaran.COD = true
	
	//dapetin cart dengan id yang sesuai dan totalkan harganya
	unwind := bson.D{{ Key: "$unwind", Value: bson.D{primitive.E{ Key: "path", Value: "$keranjang_user"}}}}
	group  := bson.D{{ Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value:"$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value:"$keranjang_user.harga"}}}}}}
	result, err := userKoleksi.Aggregate(ctx, mongo.Pipeline{unwind, group})
	ctx.Done()
	if err!=nil {
		panic(err)
	}
	var getUserCart []bson.M
	if err = result.All(ctx, &getUserCart); err != nil {
		panic(err)
	}
	var total_price int32
	for _, user := range getUserCart {
		total_price += user["total"].(int32)
	}
	orderCart.Harga = int(total_price)


	//menambahkan daftar pesanan ke user dengan id yang sesuai
	filter := bson.D{primitive.E{Key:"_id", Value: id}}
	update := bson.D{{ Key:"$push", Value: bson.D{primitive.E{Key: "pesanans", Value: orderCart}} }}
	_,err   = userKoleksi.UpdateMany(ctx, filter, update)
	if err!=nil {
		log.Println(err)
	}


	//menambahkan daftar pesanan ke koleksi pesanan dengan user id yang sesuai
	err = userKoleksi.FindOne(ctx, bson.D{primitive.E{Key:"_id", Value: id}}).Decode(&user)
	if err!=nil {
		log.Println(err)
	}
	filter2 := bson.D{{ Key:"_id", Value:id}}
	update2 := bson.D{{ Key:"$push", Value: bson.M{"pesanans.$[].list_pesanan" : bson.M{"$each": user.KeranjangUser}}}}
	_, err = userKoleksi.UpdateOne(ctx, filter2, update2)
	if err!=nil {
		log.Println(err)
	}


	//menghapus pesanan di user dengan id yang sesuai
	emptyCart := make([]models.ProdukUser, 0)
	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 :=  bson.D{{ Key: "$set", Value: bson.D{primitive.E{Key: "keranjang_user", Value: emptyCart}}}}
	_, err  = userKoleksi.UpdateOne(ctx, filter3, update3)
	if err!=nil {
		return ErrCantBuyCartItem
	}


	return nil
}

func PembeliCepat()  {
	
}