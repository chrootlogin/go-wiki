package user

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
	"encoding/json"
	"io/ioutil"
)

func TestGetUserHandler(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()

	r := gin.Default()
	r.GET("/api/user/*username", GetUserHandler)

	req, _ := http.NewRequest("GET", "/api/user/admin", nil)
	r.ServeHTTP(w, req)

	if assert.Equal(w.Code, http.StatusOK) {
		data, err := ioutil.ReadAll(w.Body)
		if assert.NoError(err) {
			var resp apiResponse
			err = json.Unmarshal(data, &resp)
			if assert.NoError(err) {
				assert.Equal("admin", resp.Username)
			}
		}
	}
}
