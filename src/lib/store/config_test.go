package store

import (
	"testing"
	"reflect"
	"github.com/chrootlogin/go-wiki/src/lib/common"
)

func TestConfig(t *testing.T) {
	config := Config()

	if !reflect.DeepEqual(common.DefaultFiles["prefs/_config.json"], config) {
		t.Error("Default configuration is not equal.")
	}
}
