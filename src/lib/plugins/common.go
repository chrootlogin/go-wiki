package plugins

import (
	"log"
	"fmt"
	"sync"
)

var instance *GoWikiPluginAPI
var once sync.Once

func GetInstance() *GoWikiPluginAPI {
	once.Do(func() {
		instance = &GoWikiPluginAPI{}
	})
	return instance
}

func RegisterPlugin(name string, pluginInterface GoWikiPlugin) {
	log.Println(fmt.Sprintf("Plugin %s (%s) was loaded!", name, pluginInterface.Version()))
}