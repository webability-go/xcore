package xcore

import (
  "fmt"
  "sync"
  "time"
  "os"
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
*/

func NewXCache(id string, maxitems int, isfile bool, expire time.Duration) *XCache {
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
}

// second boolean returned parameter is "invalid"
// if the object does not exist in memory, returns nil, false
// if the object does exist and is good, returns object, false
// if the object does exist and is invalid, returns nil, true
// Objects can become invalid when the expiration date has passed, or when the original source is newer (file type cache)
func (c *XCache)Get(key string) (interface{}, bool) {
  c.mutex.Lock()
  if x, ok := (*c).items[key]; ok {
    c.mutex.Unlock()
    if c.isfile {
      fi, err := os.Stat(key)
      if err != nil {
        // destroy the entry AND the used memory
        fmt.Println("Cache File Error and Invalid: " + key)
        return nil, true
      }
      mtime := fi.ModTime()
      if mtime.After(x.mtime) {
        // destroy the entry AND the used memory
        fmt.Println("Cache File Modified and Invalid: " + key)
        return nil, true
      }
    }
    // expired ?
        // destroy the entry AND the used memory
    
    return x.data, false
  }
  c.mutex.Unlock()
  return nil, false
}

func (c *XCache)Del(key string) {


}

func (c *XCache)Count() int {
  c.mutex.Lock()
  x := len(c.items)
  c.mutex.Unlock()
  return x
}
