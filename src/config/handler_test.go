package config

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/chrootlogin/go-wiki/src/lib/common"
	"encoding/json"
	"io/ioutil"
	"reflect"
)

func TestGetConfigHandler(t *testing.T) {
	w := httptest.NewRecorder()

	r := gin.Default()
	r.GET("/api/config", GetConfigHandler)

	req, _ := http.NewRequest("GET", "/api/config", nil)
	r.ServeHTTP(w, req)

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
	}

	var readConfiguration common.Configuration
	err = json.Unmarshal(body, &readConfiguration)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(common.DefaultFiles["prefs/_config.json"], readConfiguration) {
		t.Error("Default configuration is not equal.")
	}
}
