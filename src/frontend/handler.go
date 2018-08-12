package frontend

import (
	"time"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/chrootlogin/go-wiki/src/lib/common"
)

func GetFrontendIndexHandler(c *gin.Context) {
	content, err := Asset("index.html")
	if err != nil {
		c.JSON(http.StatusNotFound, common.ApiResponse{ Message: err.Error() })
		return
	}

	t := time.Now().AddDate(0,0,-30)
	c.Header("Expires", t.Format(time.RFC1123))
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	c.Header("Content-Type", "text/html")

	c.String(http.StatusOK, string(content))
}

func GetFrontendHandler(c *gin.Context) {
	path := c.Param("path")

	if len(path) > 0 {
		path = "assets/" + trimLeftChar(path)
	}

	content, err := Asset(path)
	if err != nil {
		c.JSON(http.StatusNotFound, common.ApiResponse{ Message: err.Error() })
		return
	}

	switch filepath.Ext(path) {
		case ".js":
			// allow caching
			t := time.Now().AddDate(0,0,365)
			c.Header("Expires", t.Format(time.RFC1123))
			c.Header("Cache-Control", "public, max-age=31536000")

			c.Header("Content-Type", "text/javascript")
		case ".html":
			// disallow caching
			t := time.Now().AddDate(0,0,-30)
			c.Header("Expires", t.Format(time.RFC1123))
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")

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