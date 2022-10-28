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

var metaCollection *mongo.Collection = config.GetCollection(config.DB, "Meta")
var validasiMeta = validator.New()

func CreateMeta(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var meta model.Meta
	defer cancel()

	if err := c.Bind(&meta) 
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status" : 400,
			"Message" : err.Error(),
		})
		return
	}

	if validationErr := validasiMeta.Struct(&meta)
	validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status" : 400,
			"Message" : validationErr.Error(),
		})
		return
	}

	newMeta := model.Meta {
		Id : primitive.NewObjectID(),
		Meta_title : meta.Meta_title,
		Meta_url : meta.Meta_url,
		Meta_desc : meta.Meta_desc,
	}

	result, err := metaCollection.InsertOne(ctx, newMeta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status" : 500,
			"Message" : err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"Status" : 200,
		"Message" : "Success create a Data!",
		"Data" : result,
	})
}

func GetMeta(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	metaId := c.Param("metaId")
	var meta model.Meta
	defer cancel()

	fmt.Println(metaId)
	objId, _ := primitive.ObjectIDFromHex(metaId)

	err := metaCollection.FindOne(ctx, bson.M{"id" : objId}).Decode(&meta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status" : 500,
			"Message" : err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status" : 200,
		"Message" : "Success get a Data",
		"Meta" : meta,
	})
}

func EditMeta(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	metaId := c.Param("metaId")
	var meta model.Meta
	defer cancel()

	fmt.Println(metaId)
	objId, _ := primitive.ObjectIDFromHex(metaId)

	if err := c.Bind(&meta)
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status" : 500,
			"Message" : err.Error(),
		})
	}

	if validationErr := validasiMeta.Struct(&meta)
	validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status" : 500,
			"Message" : validationErr.Error(),
		})
	}

	update := bson.M{"meta_title" : meta.Meta_title, "meta_desc" : meta.Meta_desc, "meta_url" : meta.Meta_url} 
	result, err := metaCollection.UpdateOne(ctx, bson.M{"id" : objId}, bson.M{"$set" : update})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status" : 400,
			"Message" : err.Error(),
		})
	}

	var updatedMeta model.Meta 
	if result.MatchedCount == 1 {
		err := metaCollection.FindOne(ctx, bson.M{"id" : objId}).Decode(&updatedMeta)

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
		"Meta" : updatedMeta,
	})
}

func DeleteMeta(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	metaId := c.Param("metaId")
	defer cancel()

	fmt.Println(metaId)
	objId, _ := primitive.ObjectIDFromHex(metaId)

	result, err := metaCollection.DeleteOne(ctx, bson.M {"id" : objId})

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

func GetAllMeta(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var meta []model.Meta
	defer cancel()

	results, err := metaCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"Status" : 500,
			"Message" : err.Error(),
		})
		return
	}

	defer results.Close(ctx) 
	for results.Next(ctx) {
		var singleMeta model.Meta
		if err = results.Decode(&singleMeta)
		err != nil {
			c.JSON(http.StatusInternalServerError, gin.H {
				"Status" : 500,
				"Message" : err.Error(),
			})
			return
		}
		meta = append(meta, singleMeta)
	}

	c.JSON(http.StatusOK, gin.H {
		"Status" : 200,
		"Message" : "Data fetched successfully!",
		"Meta" : meta,
	})
}

