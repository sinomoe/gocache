package gocache

import (
	"sync"

	"errors"
)

// Getter describes user defined function used to get data
// when data not exists in the cache
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc is a function which implements Getter interface
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group is a cache namespace and associated data loaded spread over
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var mu sync.RWMutex
var groups = make(map[string]*Group)

// NewGroup creates a new group instance and hold it in groups
func NewGroup(name string, cacheBytes int, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup return the group according to the specified group name
// if there is not a group named as name, it returns nil
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	if g, ok := groups[name]; ok {
		return g
	}
	return nil
}

// Get gets the byteview according to the specified key
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, errors.New("key is required")
	}
	// get is conccurent safe
	if v, ok := g.mainCache.get(key); ok {
		return v, nil
	}
	return g.load(key)
}

// load loads data when cache missed according to user defined get function
func (g *Group) load(key string) (ByteView, error) {
	// when cache miss, get data according to user defined get function
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	bv := ByteView{b: bytes}
	// add byteview into cache
	g.mainCache.add(key, bv)
	return bv, nil
}
