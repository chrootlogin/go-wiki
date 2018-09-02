package filemanager

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/chrootlogin/go-wiki/src/lib/common"
	"bytes"
	"mime/multipart"
	"github.com/chrootlogin/go-wiki/src/lib/filesystem"
)

func TestListFolderHandler(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()

	r := gin.Default()
	r.GET("/api/list/*path", ListFolderHandler)

	req, _ := http.NewRequest("GET", "/api/list/", nil)
	r.ServeHTTP(w, req)

	body, err := ioutil.ReadAll(w.Body)
	if assert.NoError(err) {
		var resp apiResponse
		err = json.Unmarshal(body, &resp)
		if assert.NoError(err) {
			assert.True(len(resp.Files) > 0)
		}
	}
}

func TestPostFileHandler(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()

	r := gin.Default()
	r.POST("/api/raw/*path", loginPostFileHandler)

	req, _ := http.NewRequest("POST", "/api/raw/", nil)
	r.ServeHTTP(w, req)

	body, err := ioutil.ReadAll(w.Body)
	if assert.NoError(err) {
		var resp common.ApiResponse
		err = json.Unmarshal(body, &resp)
		if assert.NoError(err) {
			assert.Equal(http.StatusInternalServerError, w.Code)
		}
	}
}

func TestPostFileHandler2(t *testing.T) {
	const TEST_CONTENT = "test content 1234"

	assert := assert.New(t)

	w := httptest.NewRecorder()

	r := gin.Default()
	r.POST("/api/raw/*path", loginPostFileHandler)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "testfile.tmp")
	if assert.NoError(err) {
		part.Write([]byte(TEST_CONTENT))

		err = writer.Close()
		if assert.NoError(err) {
			req, _ := http.NewRequest("POST", "/api/raw/", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			r.ServeHTTP(w, req)

			body, err := ioutil.ReadAll(w.Body)
			if assert.NoError(err) {
				var resp common.ApiResponse
				err = json.Unmarshal(body, &resp)
				if assert.NoError(err) {
					t.Log(resp.Message)

					assert.Equal(http.StatusCreated, w.Code)

					file, err := filesystem.New(filesystem.WithChroot("pages")).Get("testfile.tmp")
					if assert.NoError(err) {
						assert.Equal(TEST_CONTENT, file.Content)
					}
				}
			}
		}
	}
}

func loginPostFileHandler(c *gin.Context) {
	c.Set("user", common.User{
		Username: "testuser",
	})

	PostFileHandler(c)
}
