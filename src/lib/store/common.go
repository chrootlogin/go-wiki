package store

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var storeCache *cache.Cache

func init() {
	storeCache = cache.New(30*time.Minute, 10*time.Minute)
}
