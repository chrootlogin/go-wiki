package page

import (
	"os"
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
	"github.com/microcosm-cc/bluemonday"
	"github.com/imdario/mergo"
	"github.com/patrickmn/go-cache"

	"github.com/chrootlogin/go-wiki/src/lib/common"
	"github.com/chrootlogin/go-wiki/src/lib/helper"
	"github.com/chrootlogin/go-wiki/src/lib/filesystem"
)

var (
	pageCache *cache.Cache
)

type apiRequest struct {
	Content string `json:"content,omitempty"`
}

type apiResponse struct {
	Title 	     string `json:"title,omitempty"`
	Content      string `json:"content,omitempty"`
	Path         string `json:"path,omitempty"`
	LastModified string `json:"last-modified,omitempty"`
}

func init() {
	pageCache = cache.New(5*time.Minute, 10*time.Minute)
}

// CREATE
func PostPageHandler(c *gin.Context) {
	user, exists := helper.GetClientUser(c)
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

		/*err := fs.Commit(path, file, filesystem.Commit{
			Author: user,
			Message: "Created page: " + path,
		})*/
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
	path := c.Param("path")
	format := c.Query("format")
	cacheKey := path + "." + format

	// check if url is cached
	cachedResponse, found := pageCache.Get(cacheKey)
	if found {
		resp := cachedResponse.(apiResponse)

		c.Header("Last-Modified", resp.LastModified)
		c.JSON(http.StatusOK, cachedResponse)
		return
	}

	// otherwise, render
	file, path, err := getPage(path)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, common.ApiResponse{Message: "Not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
		return
	}

	lastModified := file.FileInfo.ModTime().Format(time.RFC1123)
	var response = apiResponse{
		Path: path,
		LastModified: lastModified,
	}

	switch format {
	case "no-render":
		response.Content = file.Content
	default:
		response.Content = renderPage(file.Content)
	}

	// cache page response
	pageCache.Set(cacheKey, response, cache.DefaultExpiration)

	c.Header("Last-Modified", lastModified)
	c.JSON(http.StatusOK, response)
}

// UPDATE
func PutPageHandler(c *gin.Context) {
	/*user, exists := helper.GetClientUser(c)
	if !exists {
		helper.Unauthorized(c)
		return
	}*/

	path := c.Param("path")

	fs := filesystem.New(filesystem.WithChroot("pages"), filesystem.WithMetadata())
	oldFile, err := fs.Get(path)
	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatusJSON(http.StatusNotFound, common.ApiResponse{ Message: "Not found, use POST to create." })
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
		return
	}

	var data apiRequest
	if c.BindJSON(&data) == nil {
		file := filesystem.File{
			Content: data.Content,
		}

		if err := mergo.Merge(&file, oldFile); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
			return
		}

		/*err := fs.Commit(path, file, filesystem.Commit{
			Author: user,
			Message: "Updated page: " + path,
		})*/
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
			return
		}

		pageCache.Flush()

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