package common

import (
	"path/filepath"
	"fmt"
	"plugin"
	"log"
	"sync"

	"github.com/chrootlogin/go-wiki/src/lib/event"
)

// Module is the interface implemented by types that
// register themselves as modular plug-ins.
type GoWikiModule interface {
	Init(*event.GoWikiEvents) error
	Version() string
}

type goWikiPluginRegistry struct {
	plugins map[string]GoWikiModule
}

var instance *goWikiPluginRegistry
var once sync.Once

func GetPluginRegistry() *goWikiPluginRegistry {
	once.Do(func() {
		instance = &goWikiPluginRegistry{
			plugins: make(map[string]GoWikiModule),
		}
	})
	return instance
}

func (gr goWikiPluginRegistry) Add(name string, pluginInterface GoWikiModule) {
	log.Println(fmt.Sprintf("Loading plugin: %s (%s)", name, pluginInterface.Version()))

	gr.plugins[name] = pluginInterface

	pluginInterface.Init()
}

func LoadPlugins() {
	log.Println("Starting plugins...")

	all_plugins, err := filepath.Glob("plugins/*.so")
	if err != nil {
		panic(err)
	}

	for _, filename := range (all_plugins) {
		p, err := plugin.Open(filename)
		if err != nil {
			log.Fatal(fmt.Sprintf("error: failed to load plugin: %v", err))
		}

		modObj, err := p.Lookup("Plugin")
		if err != nil {
			log.Fatal(fmt.Sprintf("error: failed to lookup type map: %v", err))
		}

		modPtr, tmapOk := modObj.(*func() interface{})
		if !tmapOk {
			log.Fatal(fmt.Sprintf("error: invalid type map: %T", modObj))
		}

		if modPtr == nil {
			log.Fatal(fmt.Sprintf("error: nil type map: type=%[1]T val=%[1]v", modPtr))
		}

		// dereference
		mod := *modPtr

		module, ok := mod().(GoWikiModule)
		if !ok {
			log.Fatal(fmt.Sprintf("error: converting: %T", modObj))
		}

		GetPluginRegistry().Add(filename, module)
	}
}