// A simple LRU cache for storing documents ([]byte). When the size maximum is reached, items are evicted
// starting with the least recently used. This data structure is goroutine-safe (it has a lock around all
// operations).
package golru_test

import (
	"github.com/manucorporat/golru"
	"math"
	"strconv"
	"testing"
)

func TestSimpleSetAndGet(t *testing.T) {
	c := golru.New(2)
	c.Set("key1", []byte("hi"))
	if string(c.Get("key1")) != "hi" || c.Len() != 1 {
		t.FailNow()
	}
}

func TestMultipleSetAndGet(t *testing.T) {
	c := golru.New(2)
	c.Set("key1", []byte("hi"))
	c.Set("key1", []byte("bye"))
	if string(c.Get("key1")) != "bye" || c.Len() != 1 {
		t.FailNow()
	}
	c.Set("key1", []byte("hey"))
	if string(c.Get("key1")) != "hey" || c.Len() != 1 {
		t.FailNow()
	}
	c.Set("key1", []byte("goodbye"))
	if string(c.Get("key1")) != "goodbye" || c.Len() != 1 {
		t.FailNow()
	}
}

func TestMultipleSetDeleteAndGet(t *testing.T) {
	c := golru.New(2)
	c.Set("key1", []byte("hi"))
	c.Del("key1")

	if c.Get("key1") != nil || c.Len() != 0 {
		t.FailNow()
	}

	if c.Get("key1") != nil || c.Len() != 0 {
		t.FailNow()
	}
}

func TestMultipleBatch(t *testing.T) {
	size := 1000
	c := golru.New(size)
	for i := 0; i < size; i++ {
		key := strconv.Itoa(i)
		c.Set(key, []byte("A"+key))
		if c.Len() != i+1 {
			t.FailNow()
		}
	}

	for i := 0; i < size; i++ {
		key := strconv.Itoa(i)
		if string(c.Get(key)) != ("A" + key) {
			t.FailNow()
		}
	}

	for i := 0; i < size; i++ {
		key := strconv.Itoa(i)
		c.Del(key)
		if c.Len() != size-i-1 {
			t.FailNow()
		}
	}
}

func TestMultipleCaped(t *testing.T) {
	c := golru.New(2)
	c.Set("key", []byte("doc0"))
	c.Set("key1", []byte("doc1"))
	c.Set("key2", []byte("doc2"))

	if c.Get("key") != nil || c.Len() != 2 {
		t.FailNow()
	}
	c.Set("key3", []byte("doc3"))
	if c.Get("key1") != nil || c.Len() != 2 {
		t.FailNow()
	}

	c.Get("key2")
	c.Set("key4", []byte("doc4"))

	if c.Get("key3") != nil || c.Len() != 2 {
		t.FailNow()
	}

	if c.Get("key2") == nil {
		t.FailNow()
	}
}

func TestMultipleCaped100(t *testing.T) {
	size := 1000
	c := golru.New(100)
	for i := 0; i < size; i++ {
		key := strconv.Itoa(i)
		c.Set(key, []byte("A"+key))
		if float64(c.Len()) != math.Min(float64(i+1), 100.0) {
			t.FailNow()
		}
	}
}

func BenchmarkSetFullSize(b *testing.B) {
	c := golru.New(b.N)
	RunSet(c, b)
}

func BenchmarkSetHalfSize(b *testing.B) {
	c := golru.New(b.N / 2)
	RunSet(c, b)
}

func BenchmarkSet10000(b *testing.B) {
	c := golru.New(10000)
	RunSet(c, b)
}

func BenchmarkGetFullSize(b *testing.B) {
	c := golru.New(b.N)
	RunGet(c, b)
}

func RunSet(c *golru.Cache, b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)
		c.Set(key, []byte(key))
	}
}

func RunGet(c *golru.Cache, b *testing.B) {
	RunSet(c, b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)
		c.Get(key)
	}
}
