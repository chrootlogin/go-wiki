package plugins

import (
	"fmt"
	"log"
	"github.com/hashicorp/go-plugin"
	"os/exec"
	"github.com/chrootlogin/go-wiki-plugin-sdk"
	"github.com/gin-gonic/gin"
)

var pluginMap = map[string]plugin.Plugin{
	"extension": &module.GoWikiPlugin{},
}

func Load(engine *gin.Engine) {
	log.Println("Starting plugins...")


	allPlugins, err := plugin.Discover("*","./plugins")
	if err != nil {
		panic(err)
	}

	for _, filename := range allPlugins {
		client := plugin.NewClient(&plugin.ClientConfig{
			Plugins: 		 pluginMap,
			HandshakeConfig: module.HandshakeConfig,
			Cmd:             exec.Command(filename),
		})
		//defer client.Kill()

		rpcClient, err := client.Client()
		if err != nil {
			log.Fatal(err)
		}

		ext, err := rpcClient.Dispense("extension")
		if err != nil {
			log.Fatal(err)
		}

		extension := ext.(module.IGoWikiPlugin)

		// Add to registry
		Registry().Add(extension.Name(), client, extension)
		log.Println(fmt.Sprintf("Plugin %s (%s) was loaded!", extension.Name(), extension.Version()))

		for _, route := range extension.Routes() {

			r := route
			var routingFunc = func(c *gin.Context) {
				var httpRequest = module.HTTPRequest{
					URL: *c.Request.URL,
				}

				c.String(200, extension.HandleRoute(r, httpRequest))
			}

			engine.GET(route, routingFunc)
		}
	}
}