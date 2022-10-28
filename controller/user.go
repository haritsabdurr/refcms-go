package controller

import (
	"context"
	"fmt"
	"golang_cms/config"
	"golang_cms/helper"
	"golang_cms/model"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var collectionUser *mongo.Collection = config.GetCollection(config.DB, "User")
var validasiUser = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([] byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([] byte(providedPassword), []byte(userPassword))
	check := true
	message := ""

	if err != nil {
		message = "Incorrect email or password!"
		check = false
	}

	return check, message
} 

func Register(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user model.User

	if err := c.BindJSON(&user)
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}

	validationErr := validasiUser.Struct(user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : validationErr.Error()})
		return
	}

	count, err := collectionUser.CountDocuments(ctx, bson.M{"email" : user.Email})
	defer cancel()
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "Error occured when checking email"})
		return
	}

	password := HashPassword(*user.Password)
	user.Password = &password

	count, err = collectionUser.CountDocuments(ctx, bson.M{"phone" : user.Phone})
	defer cancel()
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "Error occured when checking phone number"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "This email or phone number already exist!"})
		return
	}

	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()
	token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
	user.Token = &token
	user.Refresh_token = &refreshToken

	resultInsertionNumber, insertErr := collectionUser.InsertOne(ctx, user)
	if insertErr != nil {
		message := "Failed to create user account"
		c.JSON(http.StatusBadRequest, gin.H{"error" : message})
		return
	}

	defer cancel()

	c.JSON(http.StatusOK, resultInsertionNumber)
}

func Login(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user model.User
	var foundUser model.User

	if err := c.BindJSON(&user)
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}

	err := collectionUser.FindOne(ctx, bson.M{"email" : user.Email}).Decode(&foundUser)
	defer cancel()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "Incorrect email or password!"})
		return
	}

	passwordIsValid, message := VerifyPassword(*user.Password, *foundUser.Password)
	defer cancel()
	if passwordIsValid != true {
		c.JSON(http.StatusBadRequest, gin.H{"error" : message})
		return
	}

	if foundUser.Email == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "User not found!"})
		return
	}

	token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)


	helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
	err = collectionUser.FindOne(ctx, bson.M{"user_id" : foundUser.User_id}).Decode(&foundUser)


	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}

	c.JSON(http.StatusOK, foundUser)
}

func GetUsers(c *gin.Context) {
	if err := helper.CheckUserType(c, "ADMIN")
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}

	page, errs := strconv.Atoi(c.Query("page"))
	if errs != nil || page < 1 {
		page = 1
	}

	startIndex := (page - 1) * recordPerPage
	startIndex, err = strconv.Atoi(c.Query("startIndex"))

	matchStage := bson.D{{"$match", bson.D{{}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
	projectStage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"total_count", 1},
			{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
		}}}
	
	result, err := collectionUser.Aggregate(ctx, mongo.Pipeline{
		matchStage, groupStage, projectStage})
	
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "Error occured when listing user items"})
	}

	var allUsers []bson.M

	if err = result.All(ctx, &allUsers)
	err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, allUsers[0])
}

func GetUser(c *gin.Context) {
	userId := c.Param("user_id")

	if err := helper.MatchUserTypeToUid(c, userId)
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user model.User

	err := collectionUser.FindOne(ctx, bson.M{"user_id" : userId}).Decode(&user)
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUserEmail(c *gin.Context) {
	userEmail := c.Param("Email")

	if err := helper.MatchUserTypeToUid(c, userEmail)
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user model.User

	err := collectionUser.FindOne(ctx, bson.M{"Email" : user.Email}).Decode(&user)
	fmt.Println(user)
	defer cancel()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	userId := c.Param("user_id")

	if err := helper.MatchUserTypeToUid(c, userId)
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var user model.User
	
	update := bson.M{"first_name" : user.First_name, "last_name" : user.Last_name, "Password" : user.Password, "Email" : user.Email, "Phone" : user.Phone}
	result, err := collectionUser.UpdateOne(ctx, bson.M{"user_id" : userId}, bson.M{"$set" : update})
	defer cancel()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}

	var updatedUser model.User

	if result.MatchedCount == 1 {
		err := collectionUser.FindOne(ctx, bson.M{"user_id" : userId}).Decode(&updatedUser)


		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status" : 400,
				"Message" : err.Error(),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"Message" : "Success",
		"Data" : updatedUser,
	})
}
