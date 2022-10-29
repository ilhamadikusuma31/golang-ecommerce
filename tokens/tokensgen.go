package tokens

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ilhamadikusuma31/golang-ecommerce/database"
	"go.mongodb.org/mongo-driver/mongo"
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

func ValidateToken(){

}

func UpdateAllTokens(){

}