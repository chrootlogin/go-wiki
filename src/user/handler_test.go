package user

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"github.com/chrootlogin/go-wiki/src/lib/store"
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

func TestRegisterHandler(t *testing.T) {
	assert := assert.New(t)

	type Test struct {
		apiReq apiRequest
		expectedCode int
	}

	tests := []Test{
		{
			apiReq: apiRequest{
				Username: "testuser1",
				Password: "admin12345",
				Email: "test@example.org",
			},
			expectedCode: http.StatusCreated,
		},
		{
			apiReq: apiRequest{
				Username: "testuser1",
				Password: "admin12345",
				Email: "test@example.org",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			apiReq: apiRequest{
				Username: "test*ç%user2",
				Password: "admin12345",
				Email: "test@example.org",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			apiReq: apiRequest{
				Username: "testuser2",
				Password: "admin12345",
				Email: "test@exam@*ple.org",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	// enable registration
	err := store.Config().Set("registration","1")
	if assert.NoError(err) {

		for _, test := range tests {
			t.Log(test.apiReq)

			data, err := json.Marshal(test.apiReq)
			if assert.NoError(err) {
				w := httptest.NewRecorder()

				r := gin.Default()
				r.POST("/user/register", RegisterHandler)

				req, _ := http.NewRequest("POST", "/user/register", bytes.NewReader(data))
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Content-Length", string(len(data)))
				r.ServeHTTP(w, req)

				assert.Equal(test.expectedCode, w.Code)
			}
		}
	}
}