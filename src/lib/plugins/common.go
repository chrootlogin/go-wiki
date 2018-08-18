package plugins

import (
	"log"
	"fmt"
	"sync"
	"github.com/gin-gonic/gin"
)

type goWikiPluginRegistry struct {
	plugins map[string]GoWikiPlugin
}

var instance *goWikiPluginRegistry
var once sync.Once

func GetInstance() *goWikiPluginRegistry {
	once.Do(func() {
		instance = &goWikiPluginRegistry{}
	})
	return instance
}

func (gr goWikiPluginRegistry) Add(name string, pluginInterface GoWikiPlugin) {
	gr.plugins[name] = pluginInterface
}

func (gr goWikiPluginRegistry) RunEngine(engine *gin.Engine) {
	for _, v := range gr.plugins {
		v.RunEngine(engine)
	}
}

func RegisterPlugin(name string, pluginInterface GoWikiPlugin) {
	log.Println(fmt.Sprintf("Plugin %s (%s) was loaded!", name, pluginInterface.Version()))

	GetInstance().Add(name, pluginInterface)
}