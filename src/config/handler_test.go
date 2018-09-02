package config

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/chrootlogin/go-wiki/src/lib/common"
)

func TestGetConfigHandler(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()

	r := gin.Default()
	r.GET("/api/config", GetConfigHandler)

	req, _ := http.NewRequest("GET", "/api/config", nil)
	r.ServeHTTP(w, req)

	body, err := ioutil.ReadAll(w.Body)
	if assert.NoError(err) {
		var readConfiguration map[string]string
		err = json.Unmarshal(body, &readConfiguration)
		if assert.NoError(err) {
			if !reflect.DeepEqual(common.DefaultFiles["prefs/_config.json"], readConfiguration) {
				t.Error("Default configuration is not equal.")
			}
		}
	}
}
