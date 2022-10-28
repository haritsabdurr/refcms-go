package main

import (
	"golang_cms/config"
	"golang_cms/routes"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	//run database
	config.ConnectDB()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8888"
	}

	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: false,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))

	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	router.Run(":" + port)
}
