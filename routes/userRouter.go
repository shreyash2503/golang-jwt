package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shreyash2503/golang-jwt/controllers"
	"github.com/shreyash2503/golang-jwt/middleware"
)

func UserRoutes(router *gin.Engine) {
	router.Use(middleware.Authenticate)
	router.GET("/users", controllers.GetUsers)
	router.GET("/users/:id", controllers.GetUser)

}