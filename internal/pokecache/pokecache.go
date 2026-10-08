// Package pokecache provides a thread-safe, in-memory cache with time-to-live (TTL) expiration.
package pokecache

import (
	"sync"
	"time"
)

// Cache represents an in-memory key-value store with concurrent access synchronization
// and automatic entry reaping based on a configurable time-to-live interval.
type Cache struct {
	cache map[string]cacheEntry
	mux   *sync.RWMutex
}

// cacheEntry holds the raw byte payload and the timestamp at which it was cached.
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// NewCache initializes and returns a new Cache instance, starting a background goroutine
// to periodically reap entries older than the specified duration interval.
func NewCache(interval time.Duration) Cache {
	c := Cache{
		cache: make(map[string]cacheEntry),
		mux:   &sync.RWMutex{},
	}
	go c.reapLoop(interval)
	return c
}

// Add stores a raw byte slice value associated with the specified key,
// stamping the entry with the current timestamp.
func (c *Cache) Add(key string, val []byte) {
	if c == nil {
		return
	}
	if c.mux == nil {
		c.mux = &sync.RWMutex{}
	}
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.cache == nil {
		c.cache = make(map[string]cacheEntry)
	}

	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// Get retrieves the byte slice associated with key from the cache.
// It returns the value and true if found, or nil and false if absent.
func (c *Cache) Get(key string) ([]byte, bool) {
	if c == nil || c.mux == nil {
		return nil, false
	}
	// Use read lock to allow multiple goroutines to read concurrently without contention
	c.mux.RLock()
	defer c.mux.RUnlock()

	if c.cache == nil {
		return nil, false
	}

	entry, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

// reapLoop continuously listens to a ticker based on interval and triggers
// the cache reaping process on each tick.
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(interval)
	}
}

// reap removes all cache entries whose creation timestamp is older than
// the cutoff duration specified by interval.
func (c *Cache) reap(interval time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()

	cutoff := time.Now().Add(-interval)
	for k, v := range c.cache {
		if v.createdAt.Before(cutoff) {
			delete(c.cache, k)
		}
	}
}
