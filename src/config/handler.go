package config

import (
	"github.com/gin-gonic/gin"
	"github.com/chrootlogin/go-wiki/src/lib/filesystem"
	"net/http"
	"github.com/chrootlogin/go-wiki/src/lib/common"
	"encoding/json"
)

// Read
func GetConfigHandler(c *gin.Context) {
	file, err := filesystem.New(filesystem.WithChroot("prefs")).Get("_config.json")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{Message: err.Error()})
		return
	}

	var configuration common.Configuration
	err = json.Unmarshal([]byte(file.Content), &configuration)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, configuration)
}