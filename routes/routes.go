package routes

import(
	"github.com/aditikotabagi/go-ecommerce/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(Routes *gin.Engine)
{
	Routes.POST("/users/signup", controllers.SignUp())
	Routes.POST("/users/login", controllers.Login())
	Routes.POST("/admin/addproduct", controllers.ProductViewHome())
	Routes.GET("/users/productview", controllers.SearchProduct())
	Routes.GET("/users/search", controllers.SearchProductItem())
}