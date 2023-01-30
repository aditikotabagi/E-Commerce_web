package main

import(
	"github.com/aditikotabagi/go-ecommerce/controllers"
	"github.com/aditikotabagi/go-ecommerce/database"
	"github.com/aditikotabagi/go-ecommerce/middleware"
	"github.com/aditikotabagi/go-ecommerce/routes"
	"github.com/gin-gonic/gin"
)

func main(){
	port := os.Getenv("PORT")
	if port == ""{
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))
	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/Checkout", app.Checkout())
	router.GET("/InstantBuy", app.Buy())

	log.Fatal(router.Run(":" + port))
}