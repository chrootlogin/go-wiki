package pagestore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/chrootlogin/go-wiki/src/lib/filesystem"
	"github.com/chrootlogin/go-wiki/src/lib/common"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	ps := New()

	assert.NotNil(ps)
}

func TestPagestore_Commit(t *testing.T) {
	file := filesystem.File{
		Content: "This is a test file",
	}

	commit := Commit{
		Message: "Test file",
		Author: common.User{
			Username: "test",
			Email: "test@example.org",
		},
	}

	err := New().Commit("test.file", file, commit)
	if err != nil {
		t.Error(err)
	}
}

func TestPagestore_Get(t *testing.T) {
	assert := assert.New(t)

	file := filesystem.File{
		Content: "This is a second test file",
	}

	commit := Commit{
		Message: "Test file",
		Author: common.User{
			Username: "test",
			Email: "test@example.org",
		},
	}

	err := New().Commit("test2.file", file, commit)
	if assert.NoError(err) {
		f, err := New().Get("test2.file")
		if assert.NoError(err) {
			assert.Equal(file.Content, f.Content)
		}
	}
}
