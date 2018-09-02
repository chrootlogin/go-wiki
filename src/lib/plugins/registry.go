package plugins

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-plugin"

	"github.com/chrootlogin/go-wiki-plugin-sdk"
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

func (r *pluginRegistry) RegisterRoutes(engine *gin.Engine) {
	for _, ext := range r.extensions {
		for _, route := range ext.Plugin.Routes() {
			r := route
			var routingFunc = func(c *gin.Context) {
				var httpRequest = module.HTTPRequest{
					URL: c.Request.URL,
				}

				resp := ext.Plugin.HandleRoute(r, httpRequest)
				handleResponse(resp, c)
			}

			engine.GET(route, routingFunc)
		}
	}
}

func (r *pluginRegistry) Clients() map[string]*plugin.Client {
	var ret = make(map[string]*plugin.Client)

	for key, extension := range r.extensions {
		ret[key] = extension.Client
	}

	return ret
}

func handleResponse(response module.HTTPResponse, c *gin.Context) {
	for k, v := range response.Headers {
		c.Header(k, v)
	}

	c.String(response.Status, response.Body)
}
