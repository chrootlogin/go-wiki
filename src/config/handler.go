package config

import (
	"github.com/gin-gonic/gin"
		"net/http"
			"github.com/chrootlogin/go-wiki/src/lib/store"
)

// Read
func GetConfigHandler(c *gin.Context) {
	c.JSON(http.StatusOK, store.GetConfig())
}