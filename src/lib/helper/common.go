package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chrootlogin/go-wiki/src/lib/common"
	"golang.org/x/crypto/bcrypt"
)

func GetClientUser(c *gin.Context) (common.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return common.User{}, false
	}
	return user.(common.User), true
}

func Unauthorized(c *gin.Context) {
	c.Header("WWW-Authenticate", "JWT realm=gowiki")
	c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{Message: "You need to be logged in to perform this action."})
}

func Forbidden(message string, c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, common.ApiResponse{Message: message})
}

// Hash a password with bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Check a password hash
func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}