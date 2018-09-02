package store

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestUserList(t *testing.T) {
	assert := assert.New(t)

	ul := UserList()

	assert.NotNil(ul)
}
