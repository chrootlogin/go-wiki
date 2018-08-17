package common

import (
	"github.com/gin-gonic/gin"
	"path/filepath"
	"fmt"
	"plugin"
)

type GoWikiPlugin interface {
	Run(e *gin.Engine)
}

func LoadPlugins(e *gin.Engine) {
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