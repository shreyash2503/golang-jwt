package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shreyash2503/golang-jwt/database"
	"github.com/shreyash2503/golang-jwt/helpers"
	"github.com/shreyash2503/golang-jwt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func Signup(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
		return
	}


	validationErr := validate.Struct(user)
	defer cancel()
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : validationErr.Error(),
		})
		return
	}

	emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email" : user.Email})
	defer cancel()
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H {
			 "error" : "Error occured while checking for the email",
		})
		return 
	}


	phoneCount, err := userCollection.CountDocuments(ctx, bson.M{"phone" : user.Phone})
	defer cancel()
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H {
			 "error" : "Error occured while checking for the phone",
		})
		return
	}

	if emailCount > 0 || phoneCount > 0{
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Email or Phone already exists",
		})
		return

	}

	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()

	user.User_id = user.ID.Hex()

	token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
	user.Token = &token
	user.Refresh_token = &refreshToken

	resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		msg := fmt.Sprintf("User item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : msg,
		})
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, resultInsertionNumber)





}

func Login(c *gin.Context) {


}

func GetUsers(c *gin.Context){

}

func GetUser(c *gin.Context) {
	userId := c.Param("id")

	if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" :  err.Error(),
		})
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"user_id" : userId}).Decode((&user))
	defer cancel()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, user)

}