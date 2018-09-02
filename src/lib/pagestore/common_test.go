package pagestore

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chrootlogin/go-wiki/src/lib/common"
	"github.com/chrootlogin/go-wiki/src/lib/filesystem"
)

const DEFAULT_CONTENT = "This is a test!"

func TestNew(t *testing.T) {
	assert := assert.New(t)

	ps := New()

	assert.NotNil(ps)
}

func TestPagestore_Commit(t *testing.T) {
	err := commitTestFile("test.file")
	if err != nil {
		t.Error(err)
	}
}

func TestPagestore_Get(t *testing.T) {
	assert := assert.New(t)

	err := commitTestFile("test2.file")
	if assert.NoError(err) {
		f, err := New().Get("test2.file")
		if assert.NoError(err) {
			assert.Equal(DEFAULT_CONTENT, f.Content)
		}
	}
}

func TestPagestore_Has(t *testing.T) {
	assert := assert.New(t)

	err := commitTestFile("existing.file")
	if assert.NoError(err) {
		has, err := New().Has("existing.file")
		if assert.NoError(err) {
			assert.True(has)
		}
	}

	has, err := New().Has("non-existing.file")
	if assert.NoError(err) {
		assert.False(has)
	}
}

func TestPagestore_List(t *testing.T) {
	assert := assert.New(t)

	filesinfo, err := New().List("")
	if assert.NoError(err) {
		assert.True(len(filesinfo) > 0)
	}
}

func commitTestFile(path string) error {
	file := filesystem.File{
		Content: DEFAULT_CONTENT,
	}

	commit := Commit{
		Message: "Test file",
		Author: common.User{
			Username: "test",
			Email:    "test@example.org",
		},
	}

	return New().Commit(path, file, commit)
}
