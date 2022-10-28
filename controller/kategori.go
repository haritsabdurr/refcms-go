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

var kategoriCollection *mongo.Collection = config.GetCollection(config.DB, "Main Category")
var validasiKategori = validator.New()

func CreateKategori(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var kategori model.MainCategory
	defer cancel()

	if err := c.Bind(&kategori)
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"Status" : 400,
			"Message" : err.Error(),
		})
		return
	}

	newKategori := model.MainCategory {
		Id : primitive.NewObjectID(),
		Kategori_Produk: kategori.Kategori_Produk,
		Nama_produk: kategori.Nama_produk,
	}

	result, err := kategoriCollection.InsertOne(ctx, newKategori)
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

func GetKategori(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	kategoriId := c.Param("kategoriId")
	var kategori model.MainCategory
	defer cancel()

	fmt.Println(kategoriId)
	objId, _ := primitive.ObjectIDFromHex(kategoriId)

	err := kategoriCollection.FindOne(ctx, bson.M {"id" : objId}).Decode(&kategori)
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
		"Data" : kategori,
	})
}

func EditKategori(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	kategoriId := c.Param("kategoriID")
	var kategori model.MainCategory
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(kategoriId)

	if err := c.Bind(&kategori)
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"Status" : 400,
			"Message" : err.Error(),
		})
	}

	if validationErr := validasiKategori.Struct(&kategori)
	validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"Status" : 400,
			"Message" : validationErr.Error(),
		})
	}

	update := bson.M {"kategori_produk" : kategori.Kategori_Produk, "nama_produk" : kategori.Nama_produk}
	result, err := kategoriCollection.UpdateOne(ctx, bson.M {"id" : objId}, bson.M {"$set" : update})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"Status" : 400,
			"Message" : err.Error(),
		})
	}

	var updatedKategori model.MainCategory
	if result.MatchedCount == 1 {
		err := kategoriCollection.FindOne(ctx, bson.M {"id" : objId}).Decode(&updatedKategori)

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
		"Data" : updatedKategori,
	})
}

func DeleteKategori(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	kategorid := c.Param("kategoriID")
	defer cancel()

	fmt.Println(kategorid)
	objId, _ := primitive.ObjectIDFromHex(kategorid)

	result, err := kategoriCollection.DeleteOne(ctx, bson.M{"id": objId})

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

func GetAllKategori(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var kategori []model.MainCategory

	defer cancel()

	result, err := kategoriCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"Status" : 500,
			"Message" : err.Error(),
		})
		return
	}

	defer result.Close(ctx)
	for result.Next(ctx) {
		var singleKategori model.MainCategory
		if err = result.Decode(&singleKategori)
		err != nil {
			c.JSON(http.StatusInternalServerError, gin.H {
				"Status" : 500,
				"Message" : err.Error(),
			})
			return
		}

		kategori = append(kategori, singleKategori)
	}

	c.JSON(http.StatusOK, gin.H {
		"Status" : 200,
		"Message" : "Data fetced successfully!",
		"Data" : kategori,
	})
}