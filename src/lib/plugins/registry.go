package plugins

import (
	"sync"

	"github.com/chrootlogin/go-wiki-plugin-sdk"
	"github.com/hashicorp/go-plugin"
)

var registry *pluginRegistry
type pluginRegistry struct {
	extensions map[string]Extension
}

var once sync.Once

type Extension struct {
	Client *plugin.Client
	Plugin module.IGoWikiPlugin
}

func Registry() *pluginRegistry {
	once.Do(func() {
		registry = &pluginRegistry{
			extensions: make(map[string]Extension),
		}
	})

	return registry
}

func (r *pluginRegistry) Add(name string, client *plugin.Client, extension module.IGoWikiPlugin) {
	r.extensions[name] = Extension{
		Client: client,
		Plugin: extension,
	}
}

func (r *pluginRegistry) Clients() map[string]*plugin.Client {
	var ret = make(map[string]*plugin.Client)

	for key, extension := range r.extensions {
		ret[key] = extension.Client
	}

	return ret
}