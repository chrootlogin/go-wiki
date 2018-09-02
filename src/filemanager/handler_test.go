package filemanager

import (
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetConfigHandler(t *testing.T) {
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
