package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	mux   *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		cache: make(map[string]cacheEntry),
		mux:   &sync.Mutex{},
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	if c == nil {
		return
	}
	if c.mux == nil {
		c.mux = &sync.Mutex{}
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

func (c *Cache) Get(key string) ([]byte, bool) {
	if c == nil || c.mux == nil {
		return nil, false
	}
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.cache == nil {
		return nil, false
	}

	entry, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(interval)
	}
}

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
