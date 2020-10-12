package gocache

import (
	"github.com/sinomoe/gocache/lru"
	"sync"
)

// cache is a concurrent access safe lru cache
type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int
}

func (c *cache) add(key string, val ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// delay initialization
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes)
	}
	c.lru.Add(key, val)

}

func (c *cache) get(key string) (val ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
