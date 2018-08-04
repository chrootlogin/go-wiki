package page

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/chrootlogin/go-wiki/src/repo"
	"github.com/chrootlogin/go-wiki/src/common"
	"os"

	"gopkg.in/russross/blackfriday.v2"
	"github.com/microcosm-cc/bluemonday"
	)

type ApiRequest struct {
	Title   string `json:"title,omitempty"`
	Path	string `json:"path,omitempty"`
	Content string `json:"page,omitempty"`
}

type apiResponse struct {
	ContentType string `json:"contentType,omitempty"`
	Content     string `json:"page,omitempty"`
}

// READ
func GetPageHandler(c *gin.Context) {
	contentPath := c.Param("page")
	if contentPath != "" {
		contentPath = trimLeftChar(contentPath)
	}

	if contentPath == "" {
		contentPath = "_default.json"
	}

	file, err := repo.GetFile(contentPath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, common.ApiResponse{ Message: "Not found" })
			return
		}

		c.JSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
		return
	}

	if file.ContentType != "text/markdown" {
		c.JSON(http.StatusMethodNotAllowed, common.ApiResponse{ Message: "Content-type is not allowed here" })
		return
	}

	unsafe := blackfriday.Run([]byte(file.Content))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, string(html))
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}