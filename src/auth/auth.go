package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/dgrijalva/jwt-go.v3"

	"github.com/chrootlogin/go-wiki/src/lib/common"
	"github.com/chrootlogin/go-wiki/src/lib/helper"
	"github.com/chrootlogin/go-wiki/src/lib/store"
)

type AuthMiddleware struct {
	Realm            string
	Key              []byte
	SigningAlgorithm *jwt.SigningMethodHMAC
	Timeout          time.Duration
}

// Login form structure.
type ApiLogin struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func GetAuthMiddleware() *AuthMiddleware {
	signingKey := os.Getenv("SESSION_KEY")

	if len(signingKey) == 0 {
		log.Fatal("Env variable 'SESSION_KEY' must be specified")
	}

	return &AuthMiddleware{
		Realm:            "gowiki",
		Key:              []byte(signingKey),
		SigningAlgorithm: jwt.SigningMethodHS512,
		Timeout:          time.Hour * 24,
	}
}

func (am *AuthMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		loggedIn := false
		user := common.User{}

		// Get token from the Authorization header if available
		// format: Authorization: Bearer <token>
		tokenString := c.GetHeader("Authorization")
		if len(tokenString) >= 1 {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Check algorithm
				if am.SigningAlgorithm != token.Method {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return am.Key, nil
			})

			if err != nil {
				c.Header("WWW-Authenticate", "JWT realm="+am.Realm)
				c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{Message: err.Error()})
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			userId := claims["id"].(string)

			// check if user exits
			user, err = store.UserList().Get(userId)
			if err != nil {
				c.Header("WWW-Authenticate", "JWT realm="+am.Realm)
				c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{Message: err.Error()})
				return
			}

			// set login flag to logged in
			loggedIn = true
		}

		// if not logged in and trying to do a changing action
		if !loggedIn && (c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE") {
			c.Header("WWW-Authenticate", "JWT realm="+am.Realm)
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{Message: "You need to be logged in to perform this action."})
			return
		}

		c.Set("loggedIn", loggedIn)

		// on login set user identity
		if loggedIn {
			c.Set("user", user)
		}

		// go further
		c.Next()
	}
}

func (am *AuthMiddleware) LoginHandler(c *gin.Context) {
	var loginData ApiLogin

	if c.BindJSON(&loginData) != nil {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: "Missing Username or Password"})
		return
	}

	loginIsError := false

	user, err := store.UserList().Get(loginData.Username)
	if err != nil {
		loginIsError = true
		log.Println(err.Error())
	}

	if !helper.CheckPasswordHash(loginData.Password, user.PasswordHash) {
		loginIsError = true
	}

	if loginIsError {
		c.Header("WWW-Authenticate", "JWT realm="+am.Realm)
		c.JSON(http.StatusUnauthorized, common.ApiResponse{Message: "Wrong Username or Password"})
		return
	}

	// Create the token
	token := jwt.New(am.SigningAlgorithm)
	claims := token.Claims.(jwt.MapClaims)

	expire := time.Now().Add(am.Timeout)
	claims["id"] = user.Username
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = time.Now().Unix()

	tokenString, err := token.SignedString(am.Key)
	if err != nil {
		c.Header("WWW-Authenticate", "JWT realm="+am.Realm)
		c.JSON(http.StatusUnauthorized, common.ApiResponse{Message: "Creating JWT token failed!"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  tokenString,
		"expire": expire.Format(time.RFC3339),
	})
}
