package common

import (
	"fmt"
	"plugin"

	"path/filepath"
)

func LoadPlugins() {
	all_plugins, err := filepath.Glob("plugins/*.so")
	if err != nil {
		panic(err)
	}

	for _, filename := range (all_plugins) {
		fmt.Println("Filename: " + filename)

		_, err := plugin.Open(filename)
		if err != nil {
			panic(err)
		}
	}
}