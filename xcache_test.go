package xcore

import "testing"

import (
	"fmt"
)

func ExampleXCache() {
	cache := NewXCache("cacheid", 0, 0)

	fmt.Println(cache.Count())
}

func TestXCache(t *testing.T) {
	cache := NewXCache("cacheid", 0, 0)

	fmt.Println(cache.Count())
}
