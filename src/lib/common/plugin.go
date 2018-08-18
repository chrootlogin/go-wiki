package common

import (
	"path/filepath"
	"fmt"
	"plugin"
	"log"
	"github.com/chrootlogin/go-wiki/src/lib/modules"
	"sync"
	"reflect"
)

type goWikiPluginRegistry struct {
	plugins map[string]modules.Module
}

var instance *goWikiPluginRegistry
var once sync.Once

func GetPluginRegistry() *goWikiPluginRegistry {
	once.Do(func() {
		instance = &goWikiPluginRegistry{
			plugins: make(map[string]modules.Module),
		}
	})
	return instance
}

func (gr goWikiPluginRegistry) Add(name string, pluginInterface *modules.Module) {
	module := *pluginInterface

	log.Println(fmt.Sprintf("Plugin %s (%s) was loaded!", name, module.Version()))

	gr.plugins[name] = module
}
/*
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

		modObj, err := p.Lookup("GetModule")
		if err != nil {
			log.Fatal(fmt.Sprintf("error: failed to lookup module: %v\n", err))
		}

		fmt.Println(reflect.TypeOf(modObj))

		mod := modObj.(func() interface{})()
		if mod == nil {
			log.Fatal(fmt.Sprintf("error: nil module: type=%[1]T val=%[1]v\n", mod))
		}

		fmt.Println(reflect.TypeOf(mod))

		// dereference
		//mod := modPtr

		obj := mod.(modules.Module)
		/*if !ok {
			log.Fatal(fmt.Sprintf("error: cannot convert module: %T\n", mod))
		}*/

		//GetPluginRegistry().Add(filename, obj)


		fmt.Println(obj)

		/*tmapObj, err := p.Lookup("Types")
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