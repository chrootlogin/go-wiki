package store

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chrootlogin/go-wiki/src/lib/common"
)

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	config := Config()

	assert.NotNil(config)
}

func TestConfigList_GetAll(t *testing.T) {
	config := Config().GetAll()

	if !reflect.DeepEqual(common.DefaultFiles["prefs/_config.json"], config) {
		t.Error("Default configuration is not equal.")
	}
}

func TestConfigList_Get(t *testing.T) {
	assert := assert.New(t)

	title, exists := Config().Get("title")
	if assert.True(exists) {
		assert.Equal("Go-Wiki", title)
	}
}

func TestConfigList_Set(t *testing.T) {
	assert := assert.New(t)

	err := Config().Set("test-key", "blabla1234")
	if assert.NoError(err) {
		testkey, exists := Config().Get("test-key")
		if assert.True(exists) {
			assert.Equal("blabla1234", testkey)
		}
	}
}

func TestConfigList_GetDefault(t *testing.T) {
	assert := assert.New(t)

	testkey := Config().GetDefault("not-existing-test-key", "a-test-value")
	assert.Equal("a-test-value", testkey)
}
