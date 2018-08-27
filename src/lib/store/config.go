package store

import (
	"github.com/chrootlogin/go-wiki/src/lib/filesystem"
	"encoding/json"
	"github.com/chrootlogin/go-wiki/src/lib/common"
	"github.com/patrickmn/go-cache"
		"log"
	"fmt"
)

func GetConfig() (common.Configuration) {
	configI, found := storeCache.Get("config")
	if !found {
		// Load from filesystem
		file, err := filesystem.New(filesystem.WithChroot("prefs")).Get("_config.json")
		if err != nil {
			log.Fatal(fmt.Sprintf("Configuration: Cannot read file! (%s)", err))
		}

		var config common.Configuration
		err = json.Unmarshal([]byte(file.Content), &config)
		if err != nil {
			log.Fatal(fmt.Sprintf("Configuration: Cannot unmarshal json! (%s)", err))
		}

		// Write to cache
		storeCache.Set("config", config, cache.NoExpiration)

		return config
	}

	// Load from cache
	config, ok := configI.(common.Configuration)
	if !ok {
		log.Fatal("Configuration: Type conversion not possible!")
	}

	return config
}


