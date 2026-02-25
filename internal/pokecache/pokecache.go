// Package pokecache provides a time-based cache for storing API responses.
package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Cache is a thread-safe, time-based key-value store that automatically
// evicts entries older than a configured interval.
type Cache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
}

// NewCache creates a new Cache that reaps stale entries every interval.
func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries: map[string]cacheEntry{},
	}
	go cache.reapLoop(interval)
	return cache
}

// Add inserts or updates a value in the cache under the given key.
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// Get retrieves a value from the cache. The second return value reports
// whether the key was found.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.entries {
			if time.Since(entry.createdAt) > interval {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}
