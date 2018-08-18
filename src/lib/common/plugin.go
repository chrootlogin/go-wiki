package common

import (
	"path/filepath"
	"fmt"
	"plugin"
	"log"
	"github.com/chrootlogin/go-wiki/src/lib/modules"
)

/*type goWikiPluginRegistry struct {
	plugins map[string]plugins.GoWikiPlugin
}

var instance *goWikiPluginRegistry
var once sync.Once

func GetPluginRegistry() *goWikiPluginRegistry {
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
}*/

func LoadPlugins() {

	fmt.Println("Loading plugins...")
	all_plugins, err := filepath.Glob("plugins/*.so")
	if err != nil {
		panic(err)
	}

	for _, filename := range (all_plugins) {
		p, err := plugin.Open(filename)
		if err != nil {
			log.Fatal(fmt.Sprintf("error: failed to load plugin: %v\n", err))
		}

		fmt.Println(p)

		tmapObj, err := p.Lookup("Types")
		if err != nil {
			log.Fatal(fmt.Sprintf("error: failed to lookup type map: %v\n", err))
		}

		fmt.Println(tmapObj)

		// assert that the Types symbol is a *map[string]func() interface{}
		tmapPtr, tmapOk := tmapObj.(*map[string]func() interface{})
		if !tmapOk {
			log.Fatal(fmt.Sprintf("error: invalid type map: %T\n", tmapObj))
		}

		fmt.Println(tmapPtr)

		// assert that the type map pointer is not nil
		if tmapPtr == nil {
			log.Fatal(fmt.Sprintf("error: nil type map: type=%[1]T val=%[1]v\n", tmapPtr))
		}

		// dereference the type map pointer
		tmap := *tmapPtr

		fmt.Println(tmap)

		// register the plug-in's modules
		for k, v := range tmap {
			modules.RegisterModule(k, v)
		}

		/*f, err := p.Lookup("GetPlugin")
		if err != nil {
			panic(err)
		}
		fmt.Println(reflect.TypeOf(f))

		/*plug := f.(func() plugins.GoWikiPlugin)()
		fmt.Println(plug)

		GetPluginRegistry().Add(filename, plug)*/
	}
}