package user

import (
	"log"
	"fmt"
	"errors"
	"regexp"
	"net/http"

	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"

	"github.com/chrootlogin/go-wiki/src/lib/common"
	"github.com/chrootlogin/go-wiki/src/auth"
)

type ApiRequest struct {
	Name     string `json:"username"`
	Password string `json:"password"`
	EMail    string `json:"email"`
}

func RegisterHandler(c *gin.Context) {
	var data ApiRequest

	if c.BindJSON(&data) == nil {
		userList, err := auth.GetUserList()
		if err != nil {
			log.Println("Error loading users: " + err.Error())
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: "Error loading users."})
			return
		}

		if err := validateNewUser(userList, data.Name, data.Password, data.EMail); err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
			return
		}

		passwordHash, err := auth.HashPassword(data.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.ApiResponse{Message: err.Error()})
			return
		}

		user := common.User{
			Username: data.Name,
			Email: data.EMail,
			PasswordHash: passwordHash,
		}

		err = userList.Add(user)
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

func validateNewUser(userList *auth.UserList, name string, password string, email string) error {
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

	_, err = userList.Get(name)
	if err == nil {
		return errors.New("Sorry, the username '" + name + "' is already taken.")
	}

	return nil
}