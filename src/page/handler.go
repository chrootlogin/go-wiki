package page

import (
	"os"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
	"github.com/microcosm-cc/bluemonday"

	"github.com/chrootlogin/go-wiki/src/repo"
	"github.com/chrootlogin/go-wiki/src/common"
	"github.com/chrootlogin/go-wiki/src/helper"
	"github.com/chrootlogin/go-wiki/src/fs"
)

type apiRequest struct {
	Path	string `json:"path,omitempty"`
	Content string `json:"content,omitempty"`
}

type apiResponse struct {
	Title 	string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Path    string `json:"path,omitempty"`
}

// READ
func GetPageHandler(c *gin.Context) {
	path := normalizePath(c.Param("path"))
	
	data, err := fs.New(fs.WithChroot("pages")).Get(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Try index file
			path = normalizeIndexPath(c.Param("path"))

			data, err = fs.New(fs.WithChroot("pages")).Get(path)
			if err != nil {
				if os.IsNotExist(err) {
					c.JSON(http.StatusNotFound, common.ApiResponse{Message: "Not found"})
					return
				}

				c.JSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
			return
		}
	}

	format := c.Query("format")
	content := string(data)

	if format != "no-render" {
		content = renderPage(content)
	}

	c.JSON(http.StatusOK, apiResponse{
		Content: content,
		Path: path,
	})

	return
}

func PutPageHandler(c *gin.Context) {
	user, exists := common.GetClientUser(c)
	if !exists {
		helper.Unauthorized(c)
		return
	}

	path := normalizePath(c.Param("path"))

	if !repo.HasWithChroot("pages", path) {
		c.JSON(http.StatusNotFound, common.ApiResponse{ Message: "Not found, use POST to create." })
		return
	}

	var data apiRequest
	if c.BindJSON(&data) == nil {
		err := repo.SaveWithChroot("pages", path, []byte(data.Content), repo.Commit{
			Author: user,
			Message: "Updated page: " + path,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
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
	// Get user
	user, exists := common.GetClientUser(c)
	if !exists {
		helper.Unauthorized(c)
		return
	}

	// Get path
	path := normalizePath(c.Param("path"))
	if repo.HasWithChroot("pages", path) {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, common.ApiResponse{ Message: "Page already exists, use PUT to edit." })
		return
	}

	var data apiRequest
	if c.BindJSON(&data) == nil {
		dir, _ := filepath.Split(path)
		err := repo.MkdirPage(dir)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
			return
		}

		var file = common.File{
			ContentType: "text/markdown",
			Content: data.Content,
		}

		err = repo.SaveFile(path, file, repo.Commit{
			Author: user,
			Message: "Created new page: " + path,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
			return
		}

		c.JSON(http.StatusOK, common.ApiResponse{
			Message: "Created page.",
		})
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ApiResponse{Message: common.WrongAPIUsageError})
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
	if path == "/" {
		path += "index.md"
	} else {
		path += ".md"
	}

	return path
}

func normalizeIndexPath(path string) string {
	if len(path) > 0 {
		lastChar := path[len(path)-1:]

		if lastChar == "/" {
			path += "index.md"
		} else {
			path += "/index.md"
		}
	}

	return path
}