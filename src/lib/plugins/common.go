package plugins

import (
	"log"
	"fmt"
)

func RegisterPlugin(name string) {
	log.Println(fmt.Sprintf("Plugin %s was loaded!", name))
}