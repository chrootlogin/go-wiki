package common

import (
	"fmt"
	"log"
	"github.com/hashicorp/go-plugin"
	"os/exec"
	"github.com/chrootlogin/go-wiki-plugin-sdk"
)

var pluginMap = map[string]plugin.Plugin{
	"extension": &module.GoWikiPlugin{},
}

func LoadPlugins() {
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
		defer client.Kill()

		rpcClient, err := client.Client()
		if err != nil {
			log.Fatal(err)
		}

		ext, err := rpcClient.Dispense("extension")
		if err != nil {
			log.Fatal(err)
		}

		extension := ext.(module.IGoWikiPlugin)
		log.Println(fmt.Sprintf("Plugin %s (%s) was loaded!", extension.Name(), extension.Version()))
	}
}