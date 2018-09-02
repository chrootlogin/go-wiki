package page

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"encoding/json"
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/chrootlogin/go-wiki/src/lib/common"
	"github.com/chrootlogin/go-wiki/src/lib/pagestore"
)

func TestPostPageHandler(t *testing.T) {
	assert := assert.New(t)

	apiReq := apiRequest{
		Content: "This is a nice test page!",
	}

	data, err := json.Marshal(apiReq)
	if assert.NoError(err) {
		w := httptest.NewRecorder()

		r := gin.Default()
		r.POST("/api/page/*path", loginPostPageHandler)

		req, _ := http.NewRequest("POST", "/api/page/a-new-test.md", bytes.NewReader(data))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", string(len(data)))
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Error(w.Body)
		}
	}

	f, err := pagestore.New().Get("a-new-test.md")
	if assert.NoError(err) {
		assert.Equal(apiReq.Content, f.Content)
	}
}

func loginPostPageHandler(c *gin.Context) {
	c.Set("user", common.User{
		Username: "testuser",
	})

	PostPageHandler(c)
}