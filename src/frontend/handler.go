package frontend

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/chrootlogin/go-wiki/src/common"
)

func GetFrontendHandler(c *gin.Context) {
	path := c.Param("path")

	if len(path) > 0 {
		lastChar := path[len(path)-1:]

		if lastChar == "/" {
			path += "index.html"
		}

		path = trimLeftChar(path)
	} else if path == "" {
		path = "index.html"
	}

	content, err := Asset(path)
	if err != nil {
		c.JSON(http.StatusNotFound, common.ApiResponse{ Message: err.Error() })
		return
	}

	switch filepath.Ext(path) {
		case ".js":
			c.Header("Content-Type", "text/javascript")
		case ".html":
			c.Header("Content-Type", "text/html")
	}

	c.String(http.StatusOK, string(content))
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}