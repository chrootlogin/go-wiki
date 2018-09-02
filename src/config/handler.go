package config

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chrootlogin/go-wiki/src/lib/store"
)

// Read
func GetConfigHandler(c *gin.Context) {
	c.JSON(http.StatusOK, store.Config())
}