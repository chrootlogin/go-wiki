package plugins

import (
	"log"
	"fmt"
)

var API GoWikiPluginAPI

func RegisterPlugin(name string) GoWikiPluginAPIInterface {
	log.Println(fmt.Sprintf("Plugin %s was loaded!", name))

	return API
}