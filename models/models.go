package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	ID             primitive.ObjectID    `json:"_id" bson:"_id"`
	Nama_Depan	   *string               `json:"nama_depan" validate:"required,min=2,max=50"`
	Nama_Belakang  *string				 `json:"nama_belakang"`
	Password 	   *string				 `json:"password"`
	Email          *string               `json:"email" validate:"email,required"`
	Hp             *string               `json:"hp"`
	Token          *string               `json:"token"`
	Refresh_Token  *string               `json:"refresh_token"`
	Created_At     time.Time             `json:"created_at"`
	Updated_At     time.Time             `json:"updated_at"`
	User_ID  	   string			     `json:"user_id"`
	KeranjangUser  []ProdukUser          `json:"keranjang_user" bson:"keranjang_user"`
	Alamat_Detail  []Alamat				 `json:"alamat" bson:"alamat"`
	Pesanan_Status []Pesanan             `json:"pesanans" bson:"pesanans"`
}

type Produk struct{
	Produk_ID  		primitive.ObjectID   `bson:"_id"`
	Nama_Produk 	*string              `json:"nama_produk"`
	Harga 			*uint64              `json:"harga"` 
	Rating 			*uint8				 `json:"rating"`
	Gambar 		    *string              `json:"gambar"`
}

type ProdukUser struct{
	Produk_ID  		primitive.ObjectID   `bson:"_id"`
	Nama_Produk 	*string              `json:"nama_produk"`
	Harga 			int                  `json:"harga" bson:"harga"`
	Rating 			*uint8               `json:"rating"`
	Gambar 		    *string              `json:"gambar"`
}

type Alamat struct{
	Alamat_ID      primitive.ObjectID    `bson:"_id"`
	Rumah		   *string	             `json:"rumah"`
	Jalan		   *string               `json:"jalan"`
	Kota           *string               `json:"kota"`
	Kodepos        *string               `json:"kodepos"`
}


type Pesanan struct{
	Pesanan_ID          primitive.ObjectID  `bson:"_id"`
	Keranjang_Pesanan   []ProdukUser        `json:"list_pesanan" bson:"list_pesanan"`
	Orderer_At          time.Time           `json:"orderer_at" bson:"orderer_at"`
	Harga               int                 `json:"harga" bson:"harga"`
	Diskon				*int                `json:"diskon" bson:"diskon"`
	Metode_Pembayaran   Pembayaran          `json:"metode_pembayaran" bson:"metode_pembayaran"`
}

type Pembayaran struct{
	Digital bool
	COD bool
}