package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/shreyash2503/golang-jwt/helpers"
)

func Authenticate(c *gin.Context) {
	clientToken := c.Request.Header.Get("token")
	if clientToken == "" {
		c.JSON(403, gin.H{
			"error": "No token provided",
		})
		c.Abort()
		return
	}
	
	claims, err := helpers.ValidateToken(clientToken)


	c.Next()

}