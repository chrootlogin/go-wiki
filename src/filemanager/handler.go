package filemanager

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/chrootlogin/go-wiki/src/lib/common"
	"github.com/chrootlogin/go-wiki/src/lib/filesystem"
	"github.com/chrootlogin/go-wiki/src/lib/pagestore"
)

type apiResponse struct {
	Files []fileResponse `json:"files"`
	Path  string         `json:"path"`
}

type fileResponse struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	ModTime string `json:"modtime"`
	Type    string `json:"string"`
}

// READ
func ListFolderHandler(c *gin.Context) {
	path := c.Param("path")

	dirContent, err := pagestore.New().List(path)
	if err != nil {
		if err == filesystem.ErrIsFile {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.ApiResponse{Message: fmt.Sprintf("Path '%s' is not a directory!", path)})
			return
		}

		if os.IsNotExist(err) {
			c.AbortWithStatusJSON(http.StatusNotFound, common.ApiResponse{Message: fmt.Sprintf("Path '%s' not found!", path)})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{Message: err.Error()})
		return
	}

	var files []fileResponse
	for _, file := range dirContent {
		f := fileResponse{
			Name:    file.Name(),
			Size:    file.Size(),
			ModTime: file.ModTime().Format(time.RFC1123),
		}

		if file.IsDir() {
			f.Type = "directory"
		} else {
			f.Type = "file"
		}

		files = append(files, f)
	}

	c.JSON(200, apiResponse{
		Files: files,
		Path:  path,
	})
}

// Upload
func PostFileHandler(c *gin.Context) {
	/*path := c.Param("path")

	// Init Filesystem
	fs := filesystem.New(filesystem.WithChroot("pages"))

	// Get user
	user, exists := helper.GetClientUser(c)
	if !exists {
		helper.Unauthorized(c)
		return
	}*/

	form, err := c.MultipartForm()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{Message: err.Error()})
		return
	}

	fileCount := 0
	files := form.File["file"]
	for _, file := range files {
		// Open uploaded file
		f, err := file.Open()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{Message: err.Error()})
			return
		}

		// Copy buffer
		buf := bytes.NewBuffer(nil)
		_, err = io.Copy(buf, f)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{Message: err.Error()})
			return
		}

		// Close file
		f.Close()

		// Save file to folder
		/*err = fs.Commit(filepath.Join(path, file.Filename), filesystem.File{
			Content: string(buf.Bytes()),
		}, filesystem.Commit{
			Message: fmt.Sprintf("Uploaded file %s", file.Filename),
			Author: user,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{Message: err.Error()})
			return
		}*/

		fileCount++
	}

	c.JSON(http.StatusOK, common.ApiResponse{Message: fmt.Sprintf("%d files uploaded.", fileCount)})
}
