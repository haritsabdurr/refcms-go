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

var descCollection *mongo.Collection = config.GetCollection(config.DB, "Description")
var validasiDesc = validator.New()

func CreateDesc(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var desc model.Desc
	defer cancel()

	if err := c.Bind(&desc)
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"Status" : 400,
			"Message" : err.Error(),
		})
		return
	}

	if validationErr := validasiDesc.Struct(&desc)
	validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"Status" : 400,
			"Message" : validationErr.Error(),
		})
		return 
	}

	newDesc := model.Desc {
		Id: primitive.NewObjectID(),
		Title: desc.Title,
		Desc: desc.Desc,
	}

	result, err := descCollection.InsertOne(ctx, newDesc)
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

func GetDesc(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	descId := c.Param("descId")

	var desc model.Desc
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(descId)
	err := descCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&desc)

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
		"Data" : desc,
	})
}

func EditDesc(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	descId := c.Param("descId")
	var desc model.Desc
	defer cancel()

	fmt.Println(descId)
	objId, err := primitive.ObjectIDFromHex(descId)

	if err := c.Bind(&desc)
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"Status" : 400,
			"Message" : err.Error(),
		})
	}

	if validationErr := validasiDesc.Struct(&desc)
	validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"Status" : 400,
			"Message" : err.Error(),
		})
	}

	update := bson.M{"title" : desc.Title, "desc" : desc.Desc}
	result, err := descCollection.UpdateOne(ctx, bson.M{"id" : objId}, bson.M{"$set" : update})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"Status" : 400,
			"Message" : err.Error(),
		})
	}

	var updatedDesc model.Desc
	if result.MatchedCount == 1 {
		err := descCollection.FindOne(ctx, bson.M {"id" : objId}).Decode(&updatedDesc)

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
		"Data" : updatedDesc,
	})
}

func DeleteDesc(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	descId := c.Param("descId")

	defer cancel()
	fmt.Println(descId)
	objId, _ := primitive.ObjectIDFromHex(descId)

	result, err := descCollection.DeleteOne(ctx, bson.M{"id": objId})

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

func GetAllDesc(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var desc []model.Desc
	defer cancel()

	result, err := descCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"Status" : 500,
			"Message" : err.Error(),
		})
		return
	}

	defer result.Close(ctx)
	for result.Next(ctx) {
		var singleDesc model.Desc
		if err = result.Decode(&singleDesc)
		err != nil {
			c.JSON(http.StatusInternalServerError, gin.H {
				"Status" : 500,
				"Message" : err.Error(),
			})
			return
		}
		desc = append(desc, singleDesc)
	}

	c.JSON(http.StatusOK, gin.H {
		"Status" : 200,
		"Message" : "Data fetched successfully!",
		"Data" : desc,
	})
}