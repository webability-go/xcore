package xcore

import (
	"fmt"
	"testing"
	"time"
)

func ExampleNewXCache() {
	cache1 := NewXCache("cacheid1", 0, 0)
	cache2 := NewXCache("cacheid2", 100, 0)
	cache3 := NewXCache("cacheid3", 0, 100*time.Minute)

	fmt.Println("All caches empty:", cache1.Count(), cache2.Count(), cache3.Count())
	// Output: All caches empty: 0 0 0
}

func TestNewXCache(t *testing.T) {
	cache1 := NewXCache("cacheid1", 0, 0)
	cache2 := NewXCache("cacheid2", 100, 0)
	cache3 := NewXCache("cacheid3", 0, 100*time.Minute)

	t.Log("All caches empty: ", cache1.Count(), cache2.Count(), cache3.Count())
}

func TestXCache(t *testing.T) {
	cache := NewXCache("cacheid", 0, 0)
	cache.Set("id1", "Data for id1")
	cache.Set("id2", "Data for id2")
	cache.Set("id3", "Data for id3")
	t.Log("Number of elements (should be 3): ", cache.Count())
	data1, ok1 := cache.Get("id1")
	t.Log("Data for id1:", data1, ok1)
	cache.Del("id1")
	t.Log("Number of elements (should be 2): ", cache.Count())
	data2, ok2 := cache.Get("id1")
	t.Log("Data for id1: (should be nil)", data2, ok2)
}

func TestXCache_invalidated(t *testing.T) {
	cache := NewXCache("cacheid", 0, 1*time.Second)
	cache.Set("id1", "Data for id1")
	cache.Set("id2", "Data for id2")
	cache.Set("id3", "Data for id3")
	t.Log("Number of elements (should be 3): ", cache.Count())
	t.Log("Wait for 2 seconds to invalidate all 3 entries")
	time.Sleep(2 * time.Second)
	t.Log("Number of elements (should be still 3): ", cache.Count())
	data1, ok1 := cache.Get("id1")
	t.Log("Data for id1:", data1, ok1)
	t.Log("Number of elements (should be 2): ", cache.Count())
	data2, ok2 := cache.Get("id2")
	t.Log("Data for id2:", data2, ok2)
	t.Log("Number of elements (should be 1): ", cache.Count())
}
