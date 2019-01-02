package main

import (
  "fmt"
  "testing"
  "github.com/webability-go/xcore"
//  "unsafe"
)

/* TEST XCACHE */

func TestXCache(t *testing.T) {
  // Test 1: creates 100 max, no file, expires in 1 sec
  cache := xcore.NewXCache("mycache1", 100, false, 1000)
  fmt.Println(cache)

  
  
//  tr = lang.Get("")
//  if tr != "XCore" {
//    t.Errorf("NewXLanguageFromString is not working correctly")
//  }
}

