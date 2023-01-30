package controllers

import(
	"time"
	"context"
	"errors"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

	type Application struct{
		prodCollection *mongo.Collection
		userCollection *mongo.Collection

	}

	func NewApplication(prodCollection, userCollection *mongo.Collection) *Application{
		return &Application{
			prodCollection: prodCollection,
			userCollection: userCollection
		}
	}

func (app *Application) AddToCart() gin.Handler{
	return func(c *gin.Context){
		productQueryID := c.Query("id")
		if productQueryID == ""{
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product id is empty"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == ""{
			log.Println("user ID is empty")
			_, c = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err!= nil{
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeOut(context.Background(), 5*time.Second)

		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, ProductID, userQueryID)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.ItendedJSON(200, "Successfully added to cart")
	}
}

func (app *Application) RemoteItem() gin.HandlerFunc{
	return func(c *gin.Context){
		productQueryID := c.Query("id")
		if productQueryID == ""{
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product id is empty"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == ""{
			log.Println("user ID is empty")
			_, c = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err!= nil{
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeOut(context.Background(), 5*time.Second)

		defer cancel()

		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, ProductID, userQueryID)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.ItendedJSON(200, "Successfully removed item from cart")
	}

}

func GetItemFromCart() gin.Handler{
	//aggregation query: 1. match stage, 2. unwind stage, 3. group stage
	return func(c *gin.Context){
		user_id := c.Query("id")
		if user_id == ""{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error":"invalid id"})
			c.Abort()
			return
		}
		user_id, _ := primitive.ObjectIDFromHex(user_id)
		var ctx, cancel = context.WithTimeOut(context.Background(), 100*time.Second)

		defer cancel()

		var filledcart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key:"_id", Value: usert_id}}).Decode(&filledcart)
		if err!=nil{
			log.Println(err)
			c.IndentedJSON(500, "Not Found")
			return
		}
		filter_match := bson.D{{Key:"$match", Value: bson.D{primitive.E{Key:"_id", ValueL usert_id}}}}

		unwind := bson.D{{Key:"$unwind", Value: bson.D{primitive.E{Key:"path", Value:"$usercart"}}}}
		grouping := bson.D{{Key:}}
	}

}

func (app *Application) BuyFromCart() gin.Handler{
	return func(c *gin.Context){
		userQueryID := c.Query("id")
		if userQueryID == ""{
			log.Println("user id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		
		var ctx, cancel = context.WithTimeOut(context.Background(), 100*time.Second)

		defer cancel()

		err = database.BuyItemFromCart(ctx,  app.userCollection, userQueryID)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.ItendedJSON(200, "Successfully placed the order")
	}

}

func (app *Application) InstantBuy() gin.HandlerFunc{
	return func(c *gin.Context){
		productQueryID := c.Query("id")
		if productQueryID == ""{
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product id is empty"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == ""{
			log.Println("user ID is empty")
			_, c = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err!= nil{
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeOut(context.Background(), 5*time.Second)

		defer cancel()

		err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, ProductID, userQueryID)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.ItendedJSON(200, "Successfully placed the order")
	}
}
