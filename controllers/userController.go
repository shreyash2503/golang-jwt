package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shreyash2503/golang-jwt/database"
	"github.com/shreyash2503/golang-jwt/helpers"
	"github.com/shreyash2503/golang-jwt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var sessionCollection *mongo.Collection = database.OpenCollection(database.Client, "session")
var validate = validator.New()

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(hash)
}

func Signup(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User
	defer cancel()
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
	password := HashPassword(*user.Password)
	user.Password = &password

	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()

	user.User_id = user.ID.Hex()

	token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
	user.Token = &token
	user.Refresh_token = &refreshToken

	resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		msg := "User item was not created"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : msg,
		})
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, resultInsertionNumber)


}

func GetRefreshToken(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	refreshToken := c.Request.Header.Get("x-refresh-token")
	defer cancel()
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Refresh token is required",
		})
		return
	}

	accessClaims, msg := helpers.ValidateToken(refreshToken)
	defer cancel()
	if msg != "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error" : "Refresh token is not valid",

		})
		return
	}
	var user models.User

	err := userCollection.FindOne(ctx, bson.M{"user_id": accessClaims.Uid}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : "Internal Server Error",
		})
	}
	accessToken, newRefreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)

	isInserted := CreateNewSession(user, c, newRefreshToken)
	
	defer cancel()
	if !isInserted {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token" : accessToken,
		"refreshToken" : newRefreshToken,
	})
	return
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = "Email or password is incorrect"
		check = false
	}
	return check, msg
}

func Login(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User
	var foundUser models.User
	defer cancel()
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
		return

	}
	err := userCollection.FindOne(ctx, bson.M{"email" : user.Email}).Decode(&foundUser)
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : "User not found",
		})
		return
	}

	passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
	defer cancel()
	if !passwordIsValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : msg,
		})
		return
	}

	if foundUser.Email == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : "User not found",
		})
		return
	}

	token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)
	helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)

	// Check if the limit of sessions is already reached
	isInserted := CreateNewSession(user, c, refreshToken)
	defer cancel()
	if !isInserted {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : "Internal Server Error",
		})
		return
	}
	err = userCollection.FindOne(ctx, bson.M{"user_id" : foundUser.User_id}).Decode(&foundUser)
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : "User not found",

		})
		return
	}

	c.JSON(http.StatusOK, foundUser)
}

func GetUsers(c *gin.Context){
	if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
		return
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	recordPerPage, err := strconv.Atoi(c.Query(("recordPerPage")))
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}
	page, err1 := strconv.Atoi(c.Query("page"))
	if err1 != nil || page < 1 {
		page = 1
	}

	startIndex := (page - 1) * recordPerPage
	startIndex, err = strconv.Atoi(c.Query("startIndex"))
	
	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{"_id", "null"}},},
			{"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}},
		},
	},
	}

	projectStage := bson.D{
		{"$project", bson.D{
			{"_id", 0}, {"total_count", 1}, {"data", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},

		}},
	}
	result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage, groupStage, projectStage,
	})
	defer cancel()
	if err != nil { 
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : "Error occured while fetching the users",
		})
		return
	}

	var allUsers []bson.M
	if err = result.All(ctx, &allUsers); err != nil {
		log.Fatal(err)
	}
	fmt.Println(allUsers)
	c.JSON(http.StatusOK, allUsers[0])
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

func CreateNewSession(user models.User, c* gin.Context, refreshToken string) (isCreated bool) {

	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)

	sessionCount, err := sessionCollection.CountDocuments(ctx, bson.M{"user": user.ID, "valid": true})

	defer cancel()
	if err != nil {
		return false
	}

	if sessionCount > 5 {
		// Find the oldest session and update it
		// sessionCollection.FindOneAndUpdate(ctx, bson.M{"user_id" : user.User_id, "valid" : true, })
		updateOptions := options.FindOneAndUpdate().SetSort(bson.D{{"last_used", 1}})
		result := sessionCollection.FindOneAndUpdate(ctx, bson.M{"user": user.ID, "valid": true}, bson.D{{"$set", bson.D{{"valid", false}}}}, updateOptions)
		defer cancel()
		if result.Err() != nil {
			return false
		}
	}

	session := models.Session{
		ID:            primitive.NewObjectID(),
		User:          user.ID,
		Valid:         true,
		User_agent:    c.Request.UserAgent(),
		Refresh_token: refreshToken,
	}
	session.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	session.Last_used, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	_, err = sessionCollection.InsertOne(ctx, session)
	defer cancel()
	if err != nil {
		return false
	}
	return true

}