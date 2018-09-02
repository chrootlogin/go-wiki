package filesystem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	fs := New()

	assert.NotNil(fs)
}

func TestFilesystem_Get(t *testing.T) {
	assert := assert.New(t)

	f, err := New().Get("prefs/_config.json")
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(f)
}

func TestFilesystem_Save(t *testing.T) {
	assert := assert.New(t)

	file := File{
		Content: "test file",
	}

	err := New().Save("test.tmp", file)
	if err != nil {
		t.Error(err)
	}

	f, err := New().Get("test.tmp")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(file.Content, f.Content)
}

func TestFilesystem_Has(t *testing.T) {
	assert := assert.New(t)

	file := File{
		Content: "test file",
	}

	err := New().Save("test2.tmp", file)
	if assert.NoError(err) {
		has, err := New().Has("test2.tmp")
		if assert.NoError(err) {
			assert.True(has)
		}

		has, err = New().Has("no-existing-test.tmp")
		if assert.NoError(err) {
			assert.False(has)
		}
	}
}
