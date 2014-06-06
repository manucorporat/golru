// A simple LRU cache for storing documents ([]byte). When the size maximum is reached, items are evicted
// with an approximated LRU (least recently used) policy. This data structure is goroutine-safe (it has a lock around all
// operations).
package golru

import (
	"math"
	"sync"
)

const (
	DefaultLRUSamples int = 5
)

type (
	item struct {
		index uint64
		value []byte
	}

	Cache struct {
		sync.Mutex
		Capacity int
		index    uint64
		samples  int
		table    map[string]*item
	}
)

// Creates a new Cache with a maximum size of items and the number of samples used to evict the LRU entries.
// If samples <= 0, DefaultLRUSamples is used. 5 by default.
func New(capacity, samples int) *Cache {
	if samples <= 0 {
		samples = DefaultLRUSamples
	}
	return &Cache{
		Capacity: capacity,
		index:    0,
		samples:  samples,
		table:    make(map[string]*item, capacity+2),
	}
}

// Returns the total number of entries in the cache
func (c *Cache) Len() int {
	c.Lock()
	defer c.Unlock()

	return len(c.table)
}

// Insert some {key, document} into the cache. If the key already exists it would be overwritten.
func (c *Cache) Set(key string, document []byte) {
	c.Lock()
	defer func() {
		c.index++
		c.Unlock()
	}()

	c.table[key] = &item{c.index, document}
	c.trim()
}

// Removes all the entries in the cache
func (c *Cache) Flush() {
	c.Lock()
	defer c.Unlock()

	for k := range c.table {
		delete(c.table, k)
	}
	c.index = 0
}

// Get retrieves a value from the cache, it would return nil is the entry is not present
func (c *Cache) Get(key string) []byte {
	c.Lock()
	defer func() {
		c.index++
		c.Unlock()
	}()

	if elt, ok := c.table[key]; ok == true {
		elt.index = c.index
		return elt.value
	} else {
		return nil
	}
}

// Delete the document indicated by the key, if it is present.
func (c *Cache) Del(key string) {
	c.Lock()
	defer c.Unlock()

	delete(c.table, key)
}

func (c *Cache) trim() {
	toremove := len(c.table) - c.Capacity
	if toremove == 1 {
		var (
			keyToRemove string = ""
			min         uint64 = math.MaxUint64
			i           int    = 0
			iterations  int    = toremove * c.samples
		)
		for key, value := range c.table {
			if value.index < min {
				min = value.index
				keyToRemove = key
			}
			i++
			if i >= iterations {
				break
			}
		}
		if keyToRemove != "" {
			delete(c.table, keyToRemove)
		}
	}
}
