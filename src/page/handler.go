package page

import (
	"os"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/russross/blackfriday.v2"
	"github.com/microcosm-cc/bluemonday"

	"github.com/chrootlogin/go-wiki/src/repo"
	"github.com/chrootlogin/go-wiki/src/common"
)

type apiRequest struct {
	Title   string `json:"title,omitempty"`
	Path	string `json:"path,omitempty"`
	Content string `json:"page,omitempty"`
}

type apiResponse struct {
	Title 	string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

// READ
func GetPageHandler(c *gin.Context) {
	contentPath := c.Param("page")
	/*if contentPath != "" {
		contentPath = trimLeftChar(contentPath)
	}*/

	if len(contentPath) > 0 {
		lastChar := contentPath[len(contentPath)-1:]

		if lastChar == "/" {
			contentPath = "_default.json"
		}
	} else if contentPath == "" {
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

	c.JSON(http.StatusOK, apiResponse{
		Title: file.Metadata["title"],
		Content: string(html),
	})
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}