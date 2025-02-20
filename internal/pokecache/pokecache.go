package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cache    map[string]cacheEntry
	lock     sync.RWMutex
	interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	// Create global cache
	c := Cache{map[string]cacheEntry{}, sync.RWMutex{}, interval}
	// Start pruning in goroutine
	go c.readLoop()
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache[key] = cacheEntry{time.Now(), val}
	fmt.Println("Added url to cache!")
}

func (c *Cache) Get(key string) (val []byte, found bool) {
	cacheEntry, found := c.cache[key]
	fmt.Println("Returning:", found, "for url", key)
	return cacheEntry.val, found
}

func (c *Cache) readLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		<-ticker.C
		c.lock.Lock()
		for url, entry := range c.cache {
			if time.Since(entry.createdAt) > c.interval {
				delete(c.cache, url)
			}
		}
		c.lock.Unlock()
	}
}
