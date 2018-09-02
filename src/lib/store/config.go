package store

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/patrickmn/go-cache"

	"github.com/chrootlogin/go-wiki/src/lib/filesystem"
)

const CONFIGURATION_CACHE = "configuration"

type configList struct {
	configuration map[string]string
}

func Config() *configList {
	configI, found := storeCache.Get(CONFIGURATION_CACHE)
	if !found {
		// Load from filesystem
		file, err := filesystem.New(filesystem.WithChroot("prefs")).Get("_config.json")
		if err != nil {
			log.Fatal(fmt.Sprintf("Configuration: Cannot read file! (%s)", err))
		}

		var config map[string]string
		err = json.Unmarshal([]byte(file.Content), &config)
		if err != nil {
			log.Fatal(fmt.Sprintf("Configuration: Cannot unmarshal json! (%s)", err))
		}

		// Write to cache
		storeCache.Set(CONFIGURATION_CACHE, config, cache.NoExpiration)

		return &configList{
			configuration: config,
		}
	}

	// Load from cache
	config, ok := configI.(map[string]string)
	if !ok {
		log.Fatal("Configuration: Type conversion not possible!")
	}

	return &configList{
		configuration: config,
	}
}

func (cl *configList) Get(key string) (string, bool) {
	val, ok := cl.configuration[key]

	return val, ok
}

func (cl *configList) GetAll() map[string]string {
	return cl.configuration
}

func (cl *configList) GetDefault(key string, standard string) string {
	val, ok := cl.configuration[key]
	if !ok {
		return standard
	}

	return val
}

func (cl *configList) Set(key string, value string) error {
	cl.configuration[key] = value

	data, err := json.Marshal(cl.configuration)
	if err != nil {
		return err
	}

	file := filesystem.File{
		Content: string(data),
	}

	err = filesystem.New(filesystem.WithChroot("prefs")).Save("_config.json", file)
	if err != nil {
		return err
	}

	storeCache.Set(CONFIGURATION_CACHE, cl.configuration, cache.NoExpiration)
	return nil
}
