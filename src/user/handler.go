package user

import (
		"fmt"
	"errors"
	"regexp"
	"net/http"

	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"

	"github.com/chrootlogin/go-wiki/src/lib/common"
	"github.com/chrootlogin/go-wiki/src/lib/store"
	"github.com/chrootlogin/go-wiki/src/lib/helper"
)

type apiRequest struct {
	Name     string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type apiResponse struct {
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
}

func RegisterHandler(c *gin.Context) {
	// check if registration is enabled
	if !store.Config().Registration {
		helper.Forbidden("Registration is disabled!", c)
		return
	}

	var data apiRequest
	if c.BindJSON(&data) == nil {
		if err := validateNewUser(data.Name, data.Password, data.Email); err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
			return
		}

		passwordHash, err := helper.HashPassword(data.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.ApiResponse{Message: err.Error()})
			return
		}

		user := common.User{
			Username: data.Name,
			Email: data.Email,
			PasswordHash: passwordHash,
		}

		err = store.UserList().Add(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.ApiResponse{Message: "Unable to register user: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, common.ApiResponse{
			Message: fmt.Sprintf("The user %v was created.", data.Name),
		})
	} else {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: common.WrongAPIUsageError})
	}
}

func GetUserHandler(c *gin.Context) {
	userName := c.Param("username")
	if len(userName) <= 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ApiResponse{Message: common.WrongAPIUsageError})
		return
	}

	// remove first character because it's always /
	userName = trimLeftChar(userName)

	user, err := store.UserList().Get(userName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{Message: fmt.Sprintf("Can't get user list: %s", err.Error())})
		return
	}

	resp := apiResponse{
		Username: user.Username,
		Email: user.Email,
		Permissions: user.Permissions,
	}

	c.JSON(http.StatusOK, resp)
}

func validateNewUser(name string, password string, email string) error {
	if len(name) <= 3 {
		return errors.New("Username must be at least 3 chars.")
	}

	if len(name) > 20 {
		return errors.New("Username can be at most 20 chars.")
	}

	var validName = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`).MatchString
	if !validName(name) {
		return errors.New("The username can only consist of alphanumeric chars, underline and minus.")
	}

	if len(password) < 8 {
		return errors.New("Password must be at least 8 chars.")
	}

	err := checkmail.ValidateFormat(email)
	if err != nil {
		return errors.New("The email address is invalid.")
	}

	_, err = store.UserList().Get(name)
	if err == nil {
		return errors.New("Sorry, the username '" + name + "' is already taken.")
	}

	return nil
}

// https://stackoverflow.com/a/48798875
func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}