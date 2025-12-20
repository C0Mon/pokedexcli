package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	data     map[string]cacheEntry
	interval time.Duration
	mu       sync.Mutex
}

func NewCache(duration time.Duration) Cache {
	c := Cache{
		data:     map[string]cacheEntry{},
		interval: duration,
	}
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, exists := c.data[key]
	return val.val, exists
}

// Runs in the background
// Deletes cache entries that have been in too long
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.data {
			if time.Since(entry.createdAt) > c.interval {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}
