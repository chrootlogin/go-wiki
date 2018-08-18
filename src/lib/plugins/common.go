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

func RegisterPlugin(name string) GoWikiPluginAPIInterface {
	log.Println(fmt.Sprintf("Plugin %s was loaded!", name))

	return GetInstance()
}