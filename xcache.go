package xcore

import (
  "fmt"
  "sync"
  "time"
  "os"
//  "log"
)

type XCacheEntry struct {
  mtime time.Time
  data interface{}
}

type XCache struct {
  mutex sync.Mutex
  id string
  maxitems int
  isfile bool
  expire time.Duration
  items map[string]*XCacheEntry
}

/* Every cache has an ID and a flag to know if it's a cache vs a file.
   If it's a cache vs a file, then the system will check the validity vs file date
   maxitems = 0: no max limit. expire = 0: never expires
*/

func NewXCache(id string, maxitems int, isfile bool, expire time.Duration) *XCache {
  fmt.Printf("Creating cache with data {id: %s, maxitems: %d, isfile: %b, expire: %t", id, maxitems, isfile, expire)
  return &XCache{
    id: id,
    isfile: isfile,
    maxitems: maxitems,
    expire: expire,
    items: make(map[string]*XCacheEntry),
  }
}

func (c *XCache)Set(key string, indata interface{}) {
  c.mutex.Lock()
  c.items[key] = &XCacheEntry{mtime: time.Now(), data: indata}
  c.mutex.Unlock()
  if len(c.items) >= c.maxitems {
    // We need a cleaning
    go c.Clean()
  }
}

// second boolean returned parameter is "invalid"
// if the object does not exist in memory, returns nil, false
// if the object does exist and is good, returns object, false
// if the object does exist and is invalid, returns nil, true
// Objects can become invalid when the expiration date has passed, or when the original source is newer (file type cache)
func (c *XCache)Get(key string) (interface{}, bool) {
  c.mutex.Lock()
  if x, ok := (*c).items[key]; ok {
    if c.isfile {
      c.mutex.Unlock()  // We do not defer if there is disk access for speed
      fi, err := os.Stat(key)
      if err != nil {
        // destroy the entry AND the used memory
        fmt.Println("Cache File Error and Invalid: " + key)
        c.Del(key);
        return nil, true
      }
      mtime := fi.ModTime()
      if mtime.After(x.mtime) {
        // destroy the entry AND the used memory
        fmt.Println("Cache File Modified and Invalid: " + key)
        c.Del(key);
        return nil, true
      }
    }
    defer c.mutex.Unlock()
    // expired ?
    if c.expire != 0 {
      if x.mtime + c.expire < time.Now() {
        fmt.Println("Cache timeout Invalid: " + key)
        delete(c.items, key);
        return nil, true
      }
    }
    return x.data, false
  }
  c.mutex.Unlock()
  return nil, false
}

func (c *XCache)Del(key string) {
  c.mutex.Lock()
  delete(c.items, key);
  c.mutex.Unlock()
}

func (c *XCache)Count() int {
  c.mutex.Lock()
  x := len(c.items)
  c.mutex.Unlock()
  return x
}

/*
  Clean: deletes expired entries, and free 10% of max items based on time
  Returns quantity of removed entries
*/
func (c *XCache)Clean(perc int) int {
  i := 0
  c.mutex.Lock()
  // 1. clean all expired items
  for k, x := range c.items {
    if x.mtime + c.expire < time.Now() {
      fmt.Println("Cache timeout Invalid: " + k)
      delete(c.items, k)
      i++
    }
  }
  // 2. clean 10% of olders
  // **** do we consider a pile list ? unshift 10% and detroy automatically to keep *fast* the process (process is lineal)
  
  c.mutex.Unlock()
  return i
}

/*
  Clean: deletes expired entries, and free 10% of max items based on time
  Returns quantity of removed entries
*/
func (c *XCache)Verify() int {
  i := c.Clean(0)
  c.mutex.Lock()
  // 1. clean all expired items
  for k, x := range c.items {
    // verifies agains source of data
  }
  // 2. clean 10% of olders
  // **** do we consider a pile list ? unshift 10% and detroy automatically to keep *fast* the process (process is lineal)
  
  c.mutex.Unlock()
  return i
}

func (c *XCache)Flush() {
  c.mutex.Lock()
  // how to really deletes the data ? ( to free memory)
  c.items = make(map[string]*XCacheEntry)
  c.mutex.Unlock()
}
