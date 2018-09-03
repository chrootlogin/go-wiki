package store

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chrootlogin/go-wiki/src/lib/common"
)

func TestUserList(t *testing.T) {
	assert := assert.New(t)

	ul := UserList()

	assert.NotNil(ul)
}

func TestUserList_Get(t *testing.T) {
	assert := assert.New(t)

	user, err := UserList().Get("admin")
	if assert.NoError(err) {
		assert.Equal("admin", user.Username)
	}
}

func TestUserList_Get2(t *testing.T) {
	assert := assert.New(t)

	_, err := UserList().Get("no-exist")
	if assert.Error(err) {
		assert.Equal(ErrUserNotExist, err)
	}
}

func TestUserList_Add(t *testing.T) {
	assert := assert.New(t)

	newUser := common.User{
		Username: "test-user",
		Email:    "test@example.org",
	}

	err := UserList().Add(newUser)
	if assert.NoError(err) {
		user, err := UserList().Get(newUser.Username)
		if assert.NoError(err) {
			assert.Equal(newUser.Username, user.Username)
			assert.Equal(newUser.Email, user.Email)
		}
	}
}
