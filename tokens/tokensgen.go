package tokens

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ilhamadikusuma31/golang-ecommerce/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email         string
	Nama_Depan    string
	Nama_Belakang string
	Uid           string
	jwt.StandardClaims
}

var UserData *mongo.Collection = database.UserData(database.Client, "Users")
var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator(email string, namadepan string, namabelakang string, uid string) (signedtoken string, signedrefreshtoken string, err error){
	klaims := &SignedDetails{
		Email: email,
		Nama_Depan: namadepan,
		Nama_Belakang: namabelakang,
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshKlaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}


	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, klaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "","", err
	}

	refreshTokens, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshKlaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return 
	}

	return token, refreshTokens, err

}

func ValidateToken(signedtoken string) (klaims *SignedDetails, pesan string){
	token, err:= jwt.ParseWithClaims(signedtoken, &SignedDetails{}, func(token *jwt.Token)(interface{}, error){
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		pesan = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok{
		pesan = "token ga valid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix(){
		pesan = "token dah expired"
		return
	}

	return claims, pesan	
}
func UpdateAllTokens(signedtoken string, signedrefreshtoken string, userid string){

	var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
	
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key:"token", Value:signedtoken})
	updateObj = append(updateObj, bson.E{Key:"refresh_token", Value:signedrefreshtoken})
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key:"updatedat", Value: updated_at})

	upsert := true
	filter := bson.M{"user_id":userid}

	opt := options.UpdateOptions{    //kalo ga ada data, tambahkan
		Upsert : &upsert,
	}	

	_,err :=  UserData.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value:updateObj},
	}, &opt)

	defer cancel()
	if err!=nil{
		log.Panic(err)
		return 
	}

}