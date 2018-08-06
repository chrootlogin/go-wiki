package page

import (
	"os"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
	"github.com/microcosm-cc/bluemonday"

	"github.com/chrootlogin/go-wiki/src/repo"
	"github.com/chrootlogin/go-wiki/src/common"
	"path/filepath"
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
	path := normalizePath(c.Param("path"))
	
	file, err := repo.GetFile(path)
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

func PutPageHandler(c *gin.Context) {
	path := normalizePath(c.Param("path"))

	_, err := repo.GetFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, common.ApiResponse{ Message: "Not found, use POST to create." })
			return
		}

		c.JSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
		return
	}

	var data apiRequest

	if c.BindJSON(&data) == nil {
		var file = &common.File{
			ContentType: "text/markdown",
			Content: data.Content,
		}

		err = repo.SaveFile(path, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
			return
		}

		c.JSON(http.StatusOK, common.ApiResponse{
			Message: "Updated page.",
		})
	} else {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: common.WrongAPIUsageError})
	}
}

func PostPageHandler(c *gin.Context) {
	path := normalizePath(c.Param("path"))

	if repo.HasFile(path) {
		c.JSON(http.StatusMethodNotAllowed, common.ApiResponse{ Message: "Page already exists, use PUT to edit." })
		return
	}

	var data apiRequest
	if c.BindJSON(&data) == nil {

		dir, filename := filepath.Split(path)

		fmt.Println(dir)
		fmt.Println(filename)
		fmt.Println(path)

		err := repo.MkdirPage(dir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
			return
		}

		var file = &common.File{
			ContentType: "text/markdown",
			Content: data.Content,
		}

		err = repo.SaveFile(path, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
			return
		}

		c.JSON(http.StatusOK, common.ApiResponse{
			Message: "Created page.",
		})
	} else {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: common.WrongAPIUsageError})
	}
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
	output := blackfriday.Run([]byte(html), blackfriday.WithRenderer(
		blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
			AbsolutePrefix: "#/wiki",
		})))

	// Sanitize HTML
	output = bluemonday.UGCPolicy().SanitizeBytes(output)

	return string(output)
}

func normalizePath(path string) string {
	if len(path) > 0 {
		lastChar := path[len(path)-1:]

		if lastChar == "/" {
			path += "_default.json"
		} else {
			path += "/_default.json"
		}
	} else if path == "" {
		path = "_default.json"
	}

	fmt.Println(path)

	return path
}