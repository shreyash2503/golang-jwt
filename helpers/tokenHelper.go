package helpers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shreyash2503/golang-jwt/constants"
	"github.com/shreyash2503/golang-jwt/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_type  string
	jwt.RegisteredClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstName string, lastName string, userType string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails {
		Email : email,
		First_name: firstName,
		Last_name: lastName,
		Uid: uid,
		User_type: userType,
		RegisteredClaims : jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * 24)),
		},
	}

	refreshClaims := &SignedDetails {
		Uid: uid,
		RegisteredClaims : jwt.RegisteredClaims {
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * 24 * 7)),
		},

	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func UpdateAllTokens(signedToken string, signedRefereshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D
	
	//! bson.E here because primitive.E is a array of of bson.E's
	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refresh_token", signedRefereshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", Updated_at})

	upsert := true

	filter := bson.M{"user_id" : userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	//! bson.D here becuase the order of operations is important for mongoDB
	_, err := userCollection.UpdateOne(
		ctx, filter, bson.D{
			{"$set", updateObj},

		}, &opt,

	)
	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
	return
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {

			return []byte(SECRET_KEY), nil
		},

	)

	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = constants.TOKEN_INVALID
		msg = err.Error()
		return
	}

	if claims.ExpiresAt.Before(time.Now().Local()) {
		msg = constants.TOKEN_EXPIRED
		return
	}

	return claims, msg
}