package controllers

import(
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

)

var UserCollection *mongo.Collection = database.UserData(database.Client, "users") 
var productCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

func HashPassword(password string) string{
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil{
		log.Panic(err)
	}
	return string(bytes)

}

func VerifyPassword (userPassword string, givenPassword string) (bool, string){ 
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""
	if err!=nil{
		msg = "Login or Password is incorrect"
		valid = false
	}
	return valid, msg
}

func SignUp() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeOut(context.Backgroun(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BlindJSON(&user); err!=nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr :=Validate.Struct(user)
		if validationErr !=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err!=nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count>0{
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err!=nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return

		}
		if count>0{
			c.JSON(http.StatusBadRequest, gin.H{"error": "this phone no. already exists"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		token, refresh_token, _ :=generate.TokenGenerator(*user.Email, *user.FirstName, *user.LastName, user.User_ID)
		user.Token = &token
		user.RefreshToken = &refresh_token
		user.User_Cart = make([]models.ProductUser,0)
		user.Address_Details = make([]models.Address, 0)
		user.OrderStatus = make([]models.Order,0)
		_,inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "the user did not get created"})
			return
		}

		defer cancel()
		c.JSON(http.StatusCreated, "Successfully signed in")

	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeOut(context.Backgroun(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BlindJSON(&user); err!=nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		defer cancel()

		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"login or password incorrect"})
			return
		}

		PasswordIsValid, msg := VerifyPassword(*userPassword, *founduser.Password)
		defer cancel()
		if !PasswordIsValid{
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}
		token, refresh_token , _ := generate.TokenGenerator(*founduser.Email, *founduser.FirstName, *founduser.LastName, *founduser.User_ID)

		defer cancel()

		generate.UpdateAllTokens(token, refresh_token, founduser.User_ID)

		c.JSON(http.StatusFound, founduser)

	}
}

func ProductViewerAdmin() gin.HandlerFunc{

}

func SearchProduct() gin.HandlerFunc{
	return func(c *gin.Context){
		var productList []models.Product
		var ctx, cancel = context.WithTimeOut(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err!=nil{
			c.IndentedJSON(http.StatusInternalServerError, "Something went wrong. Please try after sometime")
			return
		}
		err = cursor.All(ctx, &productList)
		if err!=nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	defer cursor.Close()

	if err := cursor.err(); err!=nil{
		log.Println(err)
		c.IndentedJSON(400,"invalid")
		return
	}
	}

}

func SearchProductByQuery() gin.HandlerFunc{
	return func(c *gin.Context){
		var searchProducts []models.Product
		queryParam := c.Query("name")
		if queryParam == ""{
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid search index"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeOut(context.Background(), 100*time.Second)
		defer cancel()

		ProductCollection.Find(ctx, bson.M{"ProductName": bson.M{"$regex":queryParam}})
		if err!=nil{
			c.IndentedJSON(404, "Something went wrong while fetching the data")
			return
		}
		err = searchquerydb.All(ctx, &searchproducts)
		if err!=nil{
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchquerydb.Close(ctx)

		if err := searchquerydb.Err(); err!= nil{
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
	

		}
		defer cancel()
		c.IndentedJSON(200, searchproducts)
	}
}