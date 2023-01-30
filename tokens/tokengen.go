package tokens

import(
	"log"
	"os"
	"time"
	"github.com/aditikotabagi/E-commerce_web/Database"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/form3tech-oss/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct{
	Email string
	FirstName string
	LastName string
	Uid string
	jwt.StandardClaims
}

var UserData *mongo.Collection = database.userData(database.Client, "Users")
var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator()(email string, firstname string, lastname string, uid string)(signedtoken string, singedrefreshtoken string, err error){
	claims := &SignedDetails{
		Email: email,
		FirstName: firstname,
		LastName: lastname,
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiredAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshclaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiredAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),

		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err!=nil{
		return "", "", err
	}
	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshclaims).SignedString([]byte(SECRET_KEY))
	if err!=nil{
		log.Panic(err)
		return
	}
	return token, refreshtoken, err

}



func ValidateToken(signedtoken string)(claims *SignedDetails, msg string){
	tokens, err := jwt.ParseWithClaims(signedtoken, &SignedDetails{}, func(token *jwt.Token)(interface{},error){
		return []byte(SECRET_KEY), nil
	})
	if err!=nil{
		msg = err.Error()
		return
	}
	claims, ok := token.Claims(*SignedDetails)
	if !ok {
		msg = "The token is invalid"
		return
	}
	claims.ExpiresAt <time.Now().Local().unix(){
		msg = "Token is already expired"
		return
	}
	return claims, msg


}

func UpdateAllTokens(signedtoken string, singedrefreshtoken string, userid string ){
	var ctx, cancel = context.WithTimeOut(context.Background(), 100*time.Second)
	var updateobj primitive.D

	updateobj = append(updateobj, bson.E{Key:"token", Value: signedtoken})
	updateobj = append(updateobj, bson.E{Key:"referesh_token", Value: singedrefreshtoken})
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateobj = append(updateobj, bson.E{Key:"updatedat", Value: updated_at})

	upsert := true

	filter := bson.M{"user_id": user_id}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := UserData.UpdateOne(ctx, filter, bson.D{
		{Key:"$set", Value: updateobj},
	},
	&opt)

	defer cancel()
	if err!=nil{
		log.Panic(err)
		return
	}
}