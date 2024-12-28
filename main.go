package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shreyash2503/golang-jwt/routes"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8777"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api/v1", func(c* gin.Context) {
		c.JSON(200, gin.H{
			"success" : "Access granted",
		})
	})

	router.GET("/api/v1/health", func(c* gin.Context) {
		c.JSON(200, gin.H{
			"success" : "Server is running",
		})

	})
	
	router.Run(":" +  port)



}