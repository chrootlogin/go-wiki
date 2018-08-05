package page

import (
	"os"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
	"github.com/microcosm-cc/bluemonday"

	"github.com/chrootlogin/go-wiki/src/repo"
	"github.com/chrootlogin/go-wiki/src/common"
	"fmt"
)

type apiRequest struct {
	Path	string `json:"path,omitempty"`
	Content string `json:"content,omitempty"`
}

type apiResponse struct {
	Title 	string   `json:"title,omitempty"`
	Content string   `json:"content,omitempty"`
}

// READ
func GetPageHandler(c *gin.Context) {
	contentPath := c.Param("path")

	fmt.Println(contentPath)
	if len(contentPath) > 0 {
		lastChar := contentPath[len(contentPath)-1:]

		if lastChar == "/" {
			contentPath += "_default.json"
		}
	} else if contentPath == "" {
		contentPath = "_default.json"
	}

	fmt.Println(contentPath)

	file, err := repo.GetFile(contentPath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, common.ApiResponse{ Message: "Not found" })
			return
		}

		c.JSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
		return
	}

	format := c.Query("format")
	if format == "no-render" {
		c.JSON(http.StatusOK, apiResponse{
			Title: file.Metadata["title"],
			Content: file.Content,
		})

		return
	}

	if file.ContentType == "text/markdown" {
		c.JSON(http.StatusOK, apiResponse{
			Title: file.Metadata["title"],
			Content: renderPage(file.Content),
		})

		return
	}

	c.JSON(http.StatusMethodNotAllowed, common.ApiResponse{ Message: "Content-type is not allowed here" })
	return
}

// GET A PREVIEW
func PostPreviewHandler(c *gin.Context) {
	var data apiRequest

	if c.BindJSON(&data) == nil {
		c.JSON(http.StatusOK, apiResponse{
			Content: renderPage(data.Content),
		})
	} else {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: common.WrongAPIUsageError})
	}
}

func renderPage(html string) string {
	// Render Markdown
	output := blackfriday.Run([]byte(html))

	// Sanitize HTML
	output = bluemonday.UGCPolicy().SanitizeBytes(output)

	return string(output)
}