package cache

import (
	"sync"
	"time"
)

// GlobalCache is an in-memory cache.
var GlobalCache = NewCache()

type Cache struct {
	sync.RWMutex
	values map[string]cacheItem
}

type cacheItem struct {
	value      any
	expiration time.Time
}

func NewCache() *Cache {
	return &Cache{
		values: make(map[string]cacheItem),
	}
}

func (c *Cache) Get(key string) (any, bool) {
	c.RLock()
	defer c.RUnlock()

	item, ok := c.values[key]
	if !ok || item.expiration.Before(time.Now()) {
		return nil, false
	}
	return item.value, true
}

func (c *Cache) Set(key string, value any, duration time.Duration) {
	c.Lock()
	defer c.Unlock()

	expiration := time.Now().Add(duration)
	c.values[key] = cacheItem{
		value:      value,
		expiration: expiration,
	}
}
