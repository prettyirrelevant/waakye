package utils

import (
	"sync"
	"time"
)

// GlobalCache is a global variable for an in-memory cache.
var GlobalCache = NewCache()

// Cache represents an in-memory cache.
type Cache struct {
	sync.RWMutex
	values map[string]cacheItem
}

// cacheItem represents an item in the cache, containing the value and its expiration time.
type cacheItem struct {
	value      any
	expiration time.Time
}

// NewCache creates a new Cache instance with an empty map for the cache items.
func NewCache() *Cache {
	return &Cache{
		values: make(map[string]cacheItem),
	}
}

// Get retrieves the value associated with a key, along with a boolean indicating whether it exists and has not expired.
func (c *Cache) Get(key string) (any, bool) {
	c.RLock()
	defer c.RUnlock()

	item, ok := c.values[key]
	if !ok || item.expiration.Before(time.Now()) {
		return nil, false
	}
	return item.value, true
}

// Set adds a new key-value pair to the cache with a given expiration time
func (c *Cache) Set(key string, value any, duration time.Duration) {
	c.Lock()
	defer c.Unlock()

	expiration := time.Now().Add(duration)
	c.values[key] = cacheItem{
		value:      value,
		expiration: expiration,
	}
}
