package singleflight

import (
	"sync"
)

// call is the in-flight or finished requests
type call struct {
	val interface{}
	err error
	wg  sync.WaitGroup
}

// Group is the maintainer of diffent call
type Group struct {
	// protect all calls
	mu sync.Mutex
	m  map[string]*call
}

// Do does a key(job) only once
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	// delay initialization
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	// if the key has appeared before, wait it
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := &call{}
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
