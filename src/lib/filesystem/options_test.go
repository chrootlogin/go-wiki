package filesystem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithChroot(t *testing.T) {
	assert := assert.New(t)

	fs := New(WithChroot("prefs"))

	assert.Equal(fs.ChrootDirectory, "prefs")

	file := File{
		Content: "A test file",
	}

	err := fs.Save("testfile.tmp", file)
	if assert.NoError(err) {
		f, err := New().Get("prefs/testfile.tmp")
		if assert.NoError(err) {
			assert.Equal(file.Content, f.Content)
		}
	}
}
