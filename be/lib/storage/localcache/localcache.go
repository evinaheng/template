package localcache

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

// LocalCache to store rarely changed api result
var lc *cache.Cache

func init() {

	// default expiration 1 day
	expireTime := 24 * time.Hour
	cleanupInterval := 24 * time.Hour
	lc = cache.New(expireTime, cleanupInterval)
}

// Get local cache value
func Get(id string) (interface{}, bool) {
	return lc.Get(id)
}

// Set local cache value
func Set(id string, value interface{}, seconds int) {
	lc.Set(id, value, time.Duration(seconds)*time.Second)
}

// Delete local cache value
func Delete(id string) {
	lc.Delete(id)
}
