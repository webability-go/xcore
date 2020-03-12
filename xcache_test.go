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

	if cache1 == nil || cache1.Count() != 0 {
		t.Error("Error creating NewXCache cache1")
		return
	}
	if cache2 == nil || cache2.Count() != 0 {
		t.Error("Error creating NewXCache cache2")
		return
	}
	if cache3 == nil || cache3.Count() != 0 {
		t.Error("Error creating NewXCache cache3")
		return
	}
}

func TestXCache(t *testing.T) {
	cache := NewXCache("cacheid", 0, 0)
	cache.Set("id1", "Data for id1")
	cache.Set("id2", "Data for id2")
	cache.Set("id3", "Data for id3")
	if cache.Count() != 3 {
		t.Error("Error counting cache")
		return
	}
	data1, deleted1 := cache.Get("id1")
	if deleted1 || data1 != "Data for id1" {
		t.Error("Error getting id1 in cache")
		return
	}
	cache.Del("id1")
	if cache.Count() != 2 {
		t.Error("Error counting cache after delete")
		return
	}
	data2, deleted2 := cache.Get("id1")
	if deleted2 || data2 != nil {
		t.Error("Error getting id2 in cache")
		return
	}
}

func TestXCache_invalidated(t *testing.T) {
	cache := NewXCache("cacheid", 0, 1*time.Second)
	cache.Set("id1", "Data for id1")
	cache.Set("id2", "Data for id2")
	cache.Set("id3", "Data for id3")
	if cache.Count() != 3 {
		t.Error("Error counting cache")
		return
	}
	t.Log("Wait for 2 seconds to invalidate all 3 entries")
	time.Sleep(2 * time.Second)
	if cache.Count() != 3 {
		t.Error("Error counting cache after invalidating")
		return
	}
	data1, deleted1 := cache.Get("id1")
	if !deleted1 || data1 != nil {
		t.Error("Error getting invalidated id1 in cache")
		return
	}
	if cache.Count() != 2 {
		t.Error("Error counting cache after invalidating")
		return
	}

	data2, deleted2 := cache.Get("id2")
	if !deleted2 || data2 != nil {
		t.Error("Error getting invalidated id2 in cache")
		return
	}
	if cache.Count() != 1 {
		t.Error("Error counting cache after invalidating")
		return
	}
}

func TestXCache_cleaning(t *testing.T) {
	cache := NewXCache("cacheid", 110, 0)
	for i := 0; i < 100; i++ {
		cache.Set(fmt.Sprintf("id%d", i), i)
	}
	if cache.Count() != 100 {
		t.Error("Error counting cache")
		return
	}
	// Clean 30%
	cache.Clean(30)
	if cache.Count() != 70 {
		t.Error("Error counting cache after cleaning")
		return
	}
	cache.Flush()
	if cache.Count() != 0 {
		t.Error("Error counting cache after flushing")
		return
	}
}
