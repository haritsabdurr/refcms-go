package controller

import (
	"context"
	"fmt"
	"golang_cms/config"
	"golang_cms/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var bannerCollection *mongo.Collection = config.GetCollection(config.DB, "Banner")
var validasiBanner = validator.New()

func CreateBanner(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var banner model.Banner
	defer cancel()

	if err := c.Bind(&banner)
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"Status" : 400,
			"Message" : err.Error(),
		})
		return
	}

	if validationErr := validasiBanner.Struct(&banner)
	validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"Status" : 400,
			"Message" : validationErr.Error(),
		})
		return 
	}

	newBanner := model.Banner {
		Id: primitive.NewObjectID(),
		Banner: banner.Banner,
		Alt: banner.Alt,
		Link: banner.Link,
	}

	result, err := bannerCollection.InsertOne(ctx, newBanner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"Status" : 500,
			"Message" : err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H {
		"Status" : 200,
		"Message" : "Data created successfully!",
		"Data" : result,
	})
}

func GetBanner(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	bannerId := c.Param("bannerId")

	var banner model.Banner
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(bannerId)
	err := bannerCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&banner)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"Status" : 500,
			"Message" : err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H {
		"Status" : 200,
		"Message" : "Data fetched successfully!",
		"Data" : banner,
	})
}

func EditBanner(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	bannerId := c.Param("bannerId")
	var banner model.Banner
	defer cancel()

	fmt.Println(bannerId)
	objId, err := primitive.ObjectIDFromHex(bannerId)

	if err := c.Bind(&banner)
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"Status" : 400,
			"Message" : err.Error(),
		})
	}

	if validationErr := validasiBanner.Struct(&banner)
	validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"Status" : 400,
			"Message" : err.Error(),
		})
	}

	update := bson.M{"banner" : banner.Banner, "alt" : banner.Alt, "link" : banner.Link}
	result, err := bannerCollection.UpdateOne(ctx, bson.M{"id" : objId}, bson.M{"$set" : update})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"Status" : 400,
			"Message" : err.Error(),
		})
	}

	var updatedBanner model.Banner
	if result.MatchedCount == 1 {
		err := bannerCollection.FindOne(ctx, bson.M {"id" : objId}).Decode(&updatedBanner)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H {
				"Status" : 400,
				"Message" : err.Error(),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H {
		"Status" : 200,
		"Message" : "Data updated successfully!",
		"Data" : updatedBanner,
	})
}

func DeleteBanner(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	bannerId := c.Param("bannerId")

	defer cancel()
	fmt.Println(bannerId)
	objId, _ := primitive.ObjectIDFromHex(bannerId)

	result, err := bannerCollection.DeleteOne(ctx, bson.M{"id": objId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"Status" : 500,
			"Message" : err.Error(),
		})
	}

	if result.DeletedCount < 1 {
		c.JSON(http.StatusNotFound, gin.H {
			"Status" : 404,
			"Message" : result,
		})
	}

	c.JSON(http.StatusOK, gin.H {
		"Status" : 200,
		"Message" : "Data deleted successfully!",
	})
}

func GetAllBanner(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var banner []model.Banner
	defer cancel()

	result, err := bannerCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"Status" : 500,
			"Message" : err.Error(),
		})
		return
	}

	defer result.Close(ctx)
	for result.Next(ctx) {
		var singleBanner model.Banner
		if err = result.Decode(&singleBanner)
		err != nil {
			c.JSON(http.StatusInternalServerError, gin.H {
				"Status" : 500,
				"Message" : err.Error(),
			})
			return
		}
		banner = append(banner, singleBanner)
	}

	c.JSON(http.StatusOK, gin.H {
		"Status" : 200,
		"Message" : "Data fetched successfully!",
		"Data" : banner,
	})
}