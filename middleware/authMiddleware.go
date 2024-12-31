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
	
	accessClaims, err1 := helpers.ValidateToken(clientToken)
	if err1 != ""{
		c.JSON(403, gin.H{
			"error": "Invalid token, Please login again",
		})
	}

	if err1 != "" {
		//! Refresh token
	}
	c.Set("email", accessClaims.Email)
	c.Set("first_name", accessClaims.First_name)
	c.Set("last_name", accessClaims.Last_name)
	c.Set("uid", accessClaims.Uid)
	c.Set("user_type", accessClaims.User_type)

	c.Next()

}