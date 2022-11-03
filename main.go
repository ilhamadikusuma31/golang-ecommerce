package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ilhamadikusuma31/golang-ecommerce/controllers"
	"github.com/ilhamadikusuma31/golang-ecommerce/database"
	"github.com/ilhamadikusuma31/golang-ecommerce/middleware"
	"github.com/ilhamadikusuma31/golang-ecommerce/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	app := controllers.AplikasiBaru(database.ProdukData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Autentikasi())
	router.GET("/addtocart", app.TambahKeKeranjang())
	router.GET("/removeitem", app.HapusItem())
	router.GET("/listcart", app.DapatkanItemDariKeranjang())
	router.POST("/addaddress", controllers.TambahAlamat())
	router.PUT("/edithomeaddress", controllers.EditAlamatRumah())
	router.PUT("/editworkaddress", controllers.EditAlamatKantor())
	router.GET("/deleteaddresses", controllers.HapusAlamat())
	router.GET("/cartcheckout", app.BeliDariKeranjang())
	router.GET("/instantbuy", app.BeliCepat())
	log.Fatal(router.Run(":" + port))
}