package common

import (
	"fmt"
	"plugin"

	"path/filepath"
	"github.com/chrootlogin/go-wiki/src/lib/plugins"
	"sync"
	"log"
	"github.com/gin-gonic/gin"
)

type goWikiPluginRegistry struct {
	plugins map[string]plugins.GoWikiPlugin
}

var instance *goWikiPluginRegistry
var once sync.Once

func GetInstance() *goWikiPluginRegistry {
	once.Do(func() {
		instance = &goWikiPluginRegistry{
			plugins: make(map[string]plugins.GoWikiPlugin),
		}
	})
	return instance
}

func (gr goWikiPluginRegistry) Add(name string, pluginInterface plugins.GoWikiPlugin) {
	log.Println(fmt.Sprintf("Plugin %s (%s) was loaded!", name, pluginInterface.Version()))

	gr.plugins[name] = pluginInterface
}

func (gr goWikiPluginRegistry) RunEngine(engine *gin.Engine) {
	for _, v := range gr.plugins {
		v.RunEngine(engine)
	}
}

func LoadPlugins() {
	all_plugins, err := filepath.Glob("plugins/*.so")
	if err != nil {
		panic(err)
	}

	for _, filename := range (all_plugins) {
		fmt.Println("Filename: " + filename)

		p, err := plugin.Open(filename)
		if err != nil {
			panic(err)
		}

		f, err := p.Lookup("GetPlugin")
		if err != nil {
			panic(err)
		}

		GetInstance().Add(filename, f.(plugins.GoWikiPlugin))
	}
}