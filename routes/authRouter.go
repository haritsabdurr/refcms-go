package routes

import (
	"golang_cms/controller"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/register", controller.Register)
	incomingRoutes.POST("/users/login", controller.Login)
}