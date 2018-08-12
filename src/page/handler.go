package page

import (
	"os"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
	"github.com/microcosm-cc/bluemonday"

	"github.com/chrootlogin/go-wiki/src/repo"
	"github.com/chrootlogin/go-wiki/src/common"
	"github.com/chrootlogin/go-wiki/src/lib/helper"
	"github.com/chrootlogin/go-wiki/src/lib/filesystem"
)

type apiRequest struct {
	Content string `json:"content,omitempty"`
}

type apiResponse struct {
	Title 	string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Path    string `json:"path,omitempty"`
}

// CREATE
func PostPageHandler(c *gin.Context) {
	user, exists := common.GetClientUser(c)
	if !exists {
		helper.Unauthorized(c)
		return
	}

	path := c.Param("path")

	fs := filesystem.New(filesystem.WithChroot("pages"), filesystem.WithMetadata())
	has, err := fs.Has(path)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
		return
	}
	if has {
		c.JSON(http.StatusNotFound, common.ApiResponse{ Message: "File found, use PUT to update." })
		return
	}

	var data apiRequest
	if c.BindJSON(&data) == nil {
		file := filesystem.File{
			Content: data.Content,
			Metadata: filesystem.Metadata{
				Permissions: map[string][]string{
					// allow anyone to read
					"anyone": []string{"r"},
				},
			},
		}

		// add default permissions for author: "read, write, admin"
		file.Metadata.Permissions["u:" + user.Username] = []string{"r", "w", "a"}

		err := fs.Commit(path, file, repo.Commit{
			Author: user,
			Message: "Created page: " + path,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
			return
		}

		c.JSON(http.StatusOK, common.ApiResponse{
			Message: "Created page.",
		})
	} else {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: common.WrongAPIUsageError})
	}
}

// READ
func GetPageHandler(c *gin.Context) {
	file, path, err := getPage(c.Param("path"))
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, common.ApiResponse{Message: "Not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
		return
	}

	switch c.Query("format") {
	case "no-render":
		c.JSON(http.StatusOK, apiResponse{
			Content: file.Content,
			Path: path,
		})
	default:
		c.JSON(http.StatusOK, apiResponse{
			Content: renderPage(file.Content),
			Path: path,
		})
	}
}

// UPDATE
func PutPageHandler(c *gin.Context) {
	user, exists := common.GetClientUser(c)
	if !exists {
		helper.Unauthorized(c)
		return
	}

	path := c.Param("path")

	fs := filesystem.New(filesystem.WithChroot("pages"), filesystem.WithMetadata())
	has, err := fs.Has(path)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
		return
	}
	if !has {
		c.JSON(http.StatusNotFound, common.ApiResponse{ Message: "Not found, use POST to create." })
		return
	}

	var data apiRequest
	if c.BindJSON(&data) == nil {
		file := filesystem.File{
			Content: data.Content,
		}

		err := fs.Commit(path, file, repo.Commit{
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

func getPage(path string) (*filesystem.File, string, error) {
	fs := filesystem.New(filesystem.WithChroot("pages"), filesystem.WithMetadata())

	// possible paths
	var paths []string
	if path == "/" {
		paths = []string{
			"/index.md",
		}
	} else {
		paths = []string{
			path,
			path + ".md",
			path + "/index.md",
		}
	}

	for i := range paths {
		effectivePath := paths[i]

		file, err := fs.Get(effectivePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			} else if err == filesystem.ErrIsDir {
				continue
			} else {
				return nil, "", err
			}
		}

		return file, effectivePath, nil
	}

	return nil, "", os.ErrNotExist
}