package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type cache struct {
	data     map[string]cacheEntry
	interval time.Duration
	mu       sync.Mutex
}

func NewCache(duration time.Duration) cache {
	c := cache{
		data:     map[string]cacheEntry{},
		interval: duration,
	}
	return c
}

func (c *cache) Add(key string, val []byte) {
	c.mu.Unlock()
	c.data[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.mu.Lock()
}

func (c *cache) Get(key string) ([]byte, bool) {
	c.mu.Unlock()
	val, exists := c.data[key]
	c.mu.Lock()
	return val.val, exists
}

func (c *cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for {
		_ = ticker.C
		c.mu.Unlock()
		for key, _ := range c.data {
			if time.Since(c.data[key].createdAt) > c.interval {
				delete(c.data, key)
			}
		}
		c.mu.Lock()
	}
}
