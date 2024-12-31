package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shreyash2503/golang-jwt/controllers"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("/users/signup", controllers.Signup)
	router.POST("/users/login", controllers.Login)
	router.GET("/users/newToken", controllers.GetRefreshToken)
}