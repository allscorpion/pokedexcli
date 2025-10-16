package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	entries  map[string]cacheEntry
	lifetime time.Duration
	mux      *sync.Mutex
}

func (c *Cache) Add(key string, val []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	return entry.val, true;
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.lifetime);
	defer ticker.Stop()
	for range ticker.C {
		c.mux.Lock()
		for k, v := range c.entries {
			if time.Since(v.createdAt) > c.lifetime {
				delete(c.entries, k)
			}
		}
		c.mux.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries:  make(map[string]cacheEntry),
		lifetime: interval,
		mux:      &sync.Mutex{},
	}
	go cache.reapLoop();
	return cache;
}