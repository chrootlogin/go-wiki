package common

import (
	"github.com/gin-gonic/gin"
)

func GetClientUser(c *gin.Context) (User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return User{}, false
	}
	return user.(User), true
}