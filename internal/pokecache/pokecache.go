package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Cache for holding entries
type Cache struct {
	cache    map[string]cacheEntry
	lock     sync.RWMutex
	interval time.Duration
}

// Creates cache with specified pruning interval
func NewCache(interval time.Duration) *Cache {
	// Create global cache
	c := Cache{map[string]cacheEntry{}, sync.RWMutex{}, interval}
	// Start pruning in goroutine
	go c.readLoop()
	return &c
}

// Adds entry to cache with provided key
func (c *Cache) Add(key string, val []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache[key] = cacheEntry{time.Now(), val}
}

// Potentially retrieves entry from cache with provided key
func (c *Cache) Get(key string) (val []byte, found bool) {
	cacheEntry, found := c.cache[key]
	return cacheEntry.val, found
}

// Prunes cache based on duration interval
func (c *Cache) readLoop() {
	// Ticker gives interval for loop cycle
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		// Wait til tick
		<-ticker.C
		// Locked from writing until pruned
		c.lock.Lock()
		for url, entry := range c.cache {
			if time.Since(entry.createdAt) > c.interval {
				delete(c.cache, url)
			}
		}
		c.lock.Unlock()
	}
}
