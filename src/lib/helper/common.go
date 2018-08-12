package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chrootlogin/go-wiki/src/common"
)

func Unauthorized(c *gin.Context) {
	c.Header("WWW-Authenticate", "JWT realm=gowiki")
	c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{Message: "You need to be logged in to perform this action."})
}