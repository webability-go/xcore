package main

import (
  "fmt"
  "testing"
  "github.com/webability-go/xcore"
  "log"
  "time"
//  "unsafe"
)

/* TEST XCACHE */

func TestXCache(t *testing.T) {
  // please log
  xcore.LOG = true
  
  // Test 1: creates 100 max, no file, expires in 1 sec
  cache := xcore.NewXCache("mycache1", 100, 1000000000)   // 1s timeout
  log.Println(cache)

  // put some values
  cache.Set("one", 1);
  d200, _ := time.ParseDuration("200ms")
  time.Sleep(d200)
  cache.Set("two", 2);
  time.Sleep(d200)
  cache.Set("three", 3);
  time.Sleep(d200)
  cache.Set("four", 4);

  log.Println(cache)
  
  // wait some timeouts of 1 and 2, clean and print
  d700, _ := time.ParseDuration("700ms")
  time.Sleep(d700)
  cache.Clean(0)
  log.Println(cache)
  

  n := cache.Count()
  if n != 2 {
    t.Errorf("XCache with wrong number of elements")
  }
}

func Validate(key string, otime time.Time) bool {
  log.Println("Into validator: " + key)
  // Could be a validator against a file, a database date of modification, etc
  // In this example: invalidate always the "a" key
  return (key != "a")
}

func TestXCacheValidate(t *testing.T) {
  // please log
  xcore.LOG = true
  
  // Test 1: creates 100 max, no file, expires in 1 sec
  cache := xcore.NewXCache("mycache1", 100, 0)   // no timeout
  cache.SetValidator(Validate)

  // put some values
  cache.Set("a", 1);
  cache.Set("b", 2);
  cache.Set("c", 3);
  cache.Set("d", 4);
  log.Println(cache)

  // should invalidate the "a"
  cache.Verify()
  log.Println(cache)

  n := cache.Count()
  if n != 3 {
    t.Errorf("XCache with wrong number of elements")
  }
}

func TestXCacheCleaner(t *testing.T) {
  // please log
  xcore.LOG = true
  
  // Test 1: creates 100 max, no file, expires in 1 sec
  cache := xcore.NewXCache("mycache1", 100, 0)   // no timeout
  
  for x := 1; x < 100; x++ {
    cache.Set(fmt.Sprintf("Num%d", x), x);
  }

  // number 5 is refreshed
  cache.Get("Num5")
  
  // should launch a Clean 10% and remove 1-4, 6-11
  cache.Set("Num100", 100)
  
  // num5 should still be here
  v, _ := cache.Get("Num5")

  // should be 90 entries
  n := cache.Count()
  
  if n != 90 || v != 5 {
    t.Errorf("XCache with wrong number of elements")
  }
}





