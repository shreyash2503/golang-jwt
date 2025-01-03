package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)


func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil

	if userType != role {
		return errors.New("Unauthorized to access this resource")
	}
	return err

}

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	err = nil

	if userType == "USER" && userId != uid {
		err = errors.New("Unauthorized to access this resource")
		return err
	}

	err = CheckUserType(c, userType)
	return err
}
