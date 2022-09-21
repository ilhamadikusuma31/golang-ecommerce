package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ilhamadikusuma31/golang-ecommerce/controllers"
	"github.com/ilhamadikusuma31/golang-ecommerce/database"
	"github.com/ilhamadikusuma31/golang-ecommerce/middleware"
	"github.com/ilhamadikusuma31/golang-ecommerce/routes"
)

func main() {

	port := "8080"
	app := controllers.NewApplication(database.ProdukData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	routes.UserRoutes(router)
	routes.Use(middleware.Autentikasi())

	router.GET("/tambah-ke-keranjang", app.TambahKeKeranjang())
	router.GET("/hapus-item", app.HapusItem())
	router.GET("/keranjang-checkout", app.KeranjangCheckout())
	router.GET("/beli-cepat", app.BeliCepat())
	router.GET("/beli-cepat", app.BeliCepat())

	//run
	log.Fatal(router.Run(":"+port))
}