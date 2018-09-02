package page

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/chrootlogin/go-wiki/src/lib/common"
	"github.com/chrootlogin/go-wiki/src/lib/pagestore"
)

func TestGetPageHandler(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()

	r := gin.Default()
	r.GET("/api/page/*path", GetPageHandler)

	req, _ := http.NewRequest("GET", "/api/page/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(http.StatusOK, w.Code)
}

func TestGetPageHandler2(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()

	r := gin.Default()
	r.GET("/api/page/*path", GetPageHandler)

	req, _ := http.NewRequest("GET", "/api/page/does-not-exist/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(http.StatusNotFound, w.Code)
}

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

func TestPostPageHandler2(t *testing.T) {
	assert := assert.New(t)

	apiReq := apiRequest{
		Content: "",
	}

	data, err := json.Marshal(apiReq)
	if assert.NoError(err) {
		w := httptest.NewRecorder()

		r := gin.Default()
		r.POST("/api/page/*path", PostPageHandler)

		req, _ := http.NewRequest("POST", "/api/page/a-new-test.md", bytes.NewReader(data))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", string(len(data)))
		r.ServeHTTP(w, req)

		assert.Equal(http.StatusUnauthorized, w.Code)
	}
}

func TestPutPageHandler(t *testing.T) {
	assert := assert.New(t)

	apiReq := apiRequest{
		Content: "This is a nicer test page!",
	}

	data, err := json.Marshal(apiReq)
	if assert.NoError(err) {
		w := httptest.NewRecorder()

		r := gin.Default()
		r.PUT("/api/page/*path", loginPutPageHandler)

		req, _ := http.NewRequest("PUT", "/api/page/a-new-test.md", bytes.NewReader(data))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", string(len(data)))
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Error(w.Body)
		}
	}

	f, err := pagestore.New().Get("a-new-test.md")

	if assert.NoError(err) {
		assert.Equal(apiReq.Content, f.Content)
	}
}

func TestPutPageHandler2(t *testing.T) {
	assert := assert.New(t)

	apiReq := apiRequest{
		Content: "",
	}

	data, err := json.Marshal(apiReq)
	if assert.NoError(err) {
		w := httptest.NewRecorder()

		r := gin.Default()
		r.PUT("/api/page/*path", PutPageHandler)

		req, _ := http.NewRequest("PUT", "/api/page/a-new-test.md", bytes.NewReader(data))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", string(len(data)))
		r.ServeHTTP(w, req)

		assert.Equal(http.StatusUnauthorized, w.Code)
	}
}

func loginPostPageHandler(c *gin.Context) {
	c.Set("user", common.User{
		Username: "testuser",
	})

	PostPageHandler(c)
}

func loginPutPageHandler(c *gin.Context) {
	c.Set("user", common.User{
		Username: "testuser",
	})

	PutPageHandler(c)
}
