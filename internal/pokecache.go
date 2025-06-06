package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	data map[string]cacheEntry
	mu   sync.RWMutex
}

type cacheEntry struct {
	created time.Time
	value   []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		data: make(map[string]cacheEntry),
	}
	c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = cacheEntry{
		created: time.Now(),
		value:   value,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.data[key]
	return entry.value, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.data {
			if time.Since(entry.created) > interval {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}
