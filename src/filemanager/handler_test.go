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

func loginPostFileHandler(c *gin.Context) {
	c.Set("user", common.User{
		Username: "testuser",
	})

	PostFileHandler(c)
}