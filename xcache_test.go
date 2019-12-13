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

func TestXCache(t *testing.T) {
	cache1 := NewXCache("cacheid1", 0, 0)
	cache2 := NewXCache("cacheid2", 100, 0)
	cache3 := NewXCache("cacheid3", 0, 100*time.Minute)

	fmt.Println("All caches empty: ", cache1.Count(), cache2.Count(), cache3.Count())
}
