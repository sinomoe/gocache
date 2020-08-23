package lru

import "container/list"

// Cache is an LRU cache. not safe for concurrency!
type Cache struct {
	// MaxEntried is the maximum number of entries
	// zero means no limit
	MaxEntries int
	// OnEvicted optinonally specifies a callback function
	// to be excuted when the last entry is evicted
	OnEvicted func(key Key, val Value)
	// l is the double linkedlist storing all entries
	l *list.List
	// cache is the map to achive O(1) read time complexity
	cache map[Key]*list.Element
}

// Key must be compareble
type Key interface{}

// Value is the type of value
type Value interface{}

// entry is the value stored in list node
type entry struct {
	key Key
	val Value
}

// New returns a new Cache instance
// if maxEntries is zero, there is no limit for entries' amount
func New(maxEntries int) *Cache {
	return &Cache{
		MaxEntries: maxEntries,
		OnEvicted:  nil,
		l:          list.New(),
		cache:      make(map[Key]*list.Element),
	}
}

func (c *Cache) removeElm(el *list.Element) {
	c.l.Remove(el)
	e := el.Value.(*entry)
	delete(c.cache, e.key)
	if c.OnEvicted != nil {
		c.OnEvicted(e.key, e.val)
	}
}

// RemoveOldest removes the oldest entry from the cache
func (c *Cache) RemoveOldest() {
	if c == nil {
		return
	}
	el := c.l.Back()
	if el != nil {
		c.removeElm(el)
	}
}

// Add add key-value pair into the cache
func (c *Cache) Add(key Key, val Value) {
	if c.cache == nil {
		c.cache = make(map[Key]*list.Element)
		c.l = list.New()
	}
	// key exists
	if el, ok := c.cache[key]; ok {
		// just modify it
		c.l.MoveToFront(el)
		el.Value.(*entry).val = val
		return
	}
	// key not exists
	c.cache[key] = c.l.PushFront(&entry{key: key, val: val})
	// reached max entries numbers
	if c.MaxEntries != 0 && c.l.Len() > c.MaxEntries {
		// evite latest rarest use entry
		c.RemoveOldest()
	}
	return
}

// Remove removes the specified key if exists
func (c *Cache) Remove(key Key) {
	if c.cache == nil {
		return
	}
	if el, ok := c.cache[key]; ok {
		c.removeElm(el)
	}
}

// Get return the value according to the key
// if key not exists, ok would be false
func (c *Cache) Get(key Key) (val Value, ok bool) {
	if c.cache == nil {
		return
	}
	if el, ok := c.cache[key]; ok {
		c.l.MoveToFront(el)
		return el.Value.(*entry).val, true
	}
	return
}

// Len returns the number of entries in the cache
func (c *Cache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.l.Len()
}

// Clear purges all entries form the cache
func (c *Cache) Clear() {
	if c.cache == nil {
		return
	}
	if c.OnEvicted != nil {
		for _, el := range c.cache {
			e := el.Value.(*entry)
			c.OnEvicted(e.key, e.val)
		}
	}
	c.l = nil
	c.cache = nil
}
