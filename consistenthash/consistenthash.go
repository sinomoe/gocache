package consistenthash

import (
	"fmt"
	"hash/crc32"
	"sort"
)

// HashFn hashs bytes to uint32
type HashFn func([]byte) uint32

// Map contains all hashed keys
type Map struct {
	// hash is an editable hash function
	hash HashFn
	// replicas is the number of virtual peers of an real peer
	replicas int
	// keys are the sorted hashs of all keys
	keys []int
	// hashMap maps hash of keys to peers
	hashMap map[int]string
}

// New creates an instance of Map which stores all hashed keys
func New(replicas int, fn HashFn) *Map {
	if fn == nil {
		fn = crc32.ChecksumIEEE
	}
	return &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
}

// Add adds some keys to the hash
// a key represent a real pear
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// hash of each virtual peer
			hash := int(m.hash([]byte(fmt.Sprint(i) + key)))
			m.keys = append(m.keys, hash)
			// map replicas to real peer
			m.hashMap[hash] = key
		}
	}
	// sorted in convient to binaty searsh peers
	sort.Ints(m.keys)
}

// Get gets the closest item in the hash to the provided key
func (m *Map) Get(keys string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(keys)))
	// binary search appropriate replica peer
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
