package database

import (
	"errors"
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


func TambahProdukKeKeranjang()  {
	
}

func HapusProdukDariKeranjang()  {
	

}
func BeliProdukDariKeranjang()  {
	
}

func PembeliCepat()  {
	
}