package xcore

import (
	"log"
	"sync"
	"time"
)

/* XCacheEntry:
  ------------
  The cache entry has a time to measure expiration if needed, or time of entry in cache.
  - ctime is creation time (used to validate the object against its source).
  - rtime is last read time (used to clean the cache: the less accessed objects are removed).
  The data as itself is an interface to whatever the user need to cache.
*/
type XCacheEntry struct {
	ctime time.Time
	rtime time.Time
	data  interface{}
}

/* XCache:
  ------
  The XCache has an id (informative).
  - The user can creates a cache with a maximum number of elements if need. In this case, when the cache reaches the maximum number of elements stored, then the system makes a clean of 10% of oldest elements. This type of use is not recommended since is it heavy in CPU use to clean the cache.
  - The user can also create an expiration duration, so every elements in the cache is invalidated after a certain amount of time. It is more recommended to use the cache with an expiration duration. The obsolete objects are destroyed when the user tries to use them and return a "non existence" on Get. (this does not use CPU or extra locks.
  - The Validator is a function that can be set to check the validity of the data (for instance if the data originates from a file or a database). The validator is called for each Get (and can be heavy for CPU or can wait a long time, for instance if the check is an external database on another cluster). Beware of this.
  - The cache owns a mutex to lock access to data to read/write/delete/clean the data, to allow concurrency and multithreading of the cache.
  - The pile keeps the "ordered by date of reading" object keys, so it's fast to clean the data.
  - Finally, the items are a map to cache entries, acceved by the key of entries.
*/
type XCache struct {
	mutex     sync.Mutex
	id        string
	maxitems  int
	validator func(string, time.Time) bool
	expire    time.Duration
	items     map[string]*XCacheEntry
	pile      []string
}

/* NewCache:
  ---------
  Creates a new XCache structure.
  The XCache is resident in memory, supports multithreading and concurrency.
  "id" is the unique id of the XCache.
  maxitems is the max authorized quantity of objects into the XCache. If 0, no limit
  expire is a max duration of the objects into the cache. If 0, no limit
  Returns the *XCache created.
*/
func NewXCache(id string, maxitems int, expire time.Duration) *XCache {
	if LOG {
		log.Printf("Creating cache with data {id: %s, maxitems: %d, expire: %d}", id, maxitems, expire)
	}
	return &XCache{
		id:        id,
		maxitems:  maxitems,
		validator: nil,
		expire:    expire,
		items:     make(map[string]*XCacheEntry),
	}
}

/* GetId:
  ---------
  exposes the ID of the cache
*/
func (c *XCache) GetId() string {
	return c.id
}

/* GetMax:
  ---------
  exposes the max quantity of items of the cache
*/
func (c *XCache) GetMax() int {
	return c.maxitems
}

/* GetExpire:
  ---------
  exposes the expiration time of the cache
*/
func (c *XCache) GetExpire() time.Duration {
	return c.expire
}

/* SetValidator:
  -------------
  Sets the validator function to check every entry in the cache against its original source.
  Returns nothing.
*/
func (c *XCache) SetValidator(f func(string, time.Time) bool) {
	c.validator = f
}

/* Set:
  ----
  Sets an entry in the cache.
  If the entry already exists, just replace it with a new creation date.
  If the entry does not exist, it will insert it in the cache and if the cache if full (maxitems reached), then a clean is called to remove 10%.
  Returns nothing.
*/
func (c *XCache) Set(key string, indata interface{}) {
	c.mutex.Lock()
	// check if the entry already exists
	_, ok := c.items[key]
	c.items[key] = &XCacheEntry{ctime: time.Now(), rtime: time.Now(), data: indata}
	if ok {
		c.removeFromPile(key)
	}
	c.pile = append(c.pile, key)
	c.mutex.Unlock()
	if c.maxitems > 0 && len(c.items) >= c.maxitems {
		// We need a cleaning
		c.Clean(10)
	}
}

/* removeFromPile:
  ---------------
  will remove an entry key from the ordered pile.
*/
func (c *XCache) removeFromPile(key string) {
	// removes the key and append it to the end
	for i, x := range c.pile {
		if x == key {
			if i == len(c.pile)-1 {
				c.pile = c.pile[:i]
			} else {
				c.pile = append(c.pile[:i], c.pile[i+1:]...)
			}
			break
		}
	}
}

/* Get:
  ----
  get the value of an entry.
  If the entry does not exists, returns nil, false.
  If the entry exists and is invalidated by time or validator function, then returns nil, true.
  If the entry is good, return <value>, false.
*/
func (c *XCache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	if x, ok := c.items[key]; ok {
		c.mutex.Unlock()
		if c.validator != nil {
			if b := c.validator(key, x.ctime); !b {
				if LOG {
					log.Println("Validator invalids entry: " + key)
				}
				c.mutex.Lock()
				delete(c.items, key)
				c.removeFromPile(key)
				c.mutex.Unlock()
				return nil, true
			}
		}
		// expired ?
		if c.expire != 0 {
			if x.ctime.Add(c.expire).Before(time.Now()) {
				if LOG {
					log.Println("Cache timeout Expired: " + key)
				}
				c.mutex.Lock()
				delete(c.items, key)
				c.removeFromPile(key)
				c.mutex.Unlock()
				return nil, true
			}
		}
		x.rtime = time.Now()
		c.removeFromPile(key)
		c.pile = append(c.pile, key)
		return x.data, false
	}
	c.mutex.Unlock()
	return nil, false
}

/* Del:
  ----
  deletes the entry of the cache if it exists.
*/
func (c *XCache) Del(key string) {
	c.mutex.Lock()
	delete(c.items, key)
	// we should check if the entry exists before trying to removing
	c.removeFromPile(key)
	c.mutex.Unlock()
}

/* Count:
  ------
  returns the quantity of entries in the cache.
*/
func (c *XCache) Count() int {
	c.mutex.Lock()
	x := len(c.items)
	c.mutex.Unlock()
	return x
}

/* Clean:
  ------
  deletes expired entries, and free perc% of max items based on time.
  perc = 0 to 100 (percentage to clean).
  Returns quantity of removed entries.
  It Will **not** verify the cache against its source (if Validator is set). If you want to scan that, use the Verify function.
*/
func (c *XCache) Clean(perc int) int {
	if LOG {
		log.Println("Cleaning cache")
	}
	i := 0
	c.mutex.Lock()
	// 1. clean all expired items
	if c.expire != 0 {
		for k, x := range c.items {
			if x.ctime.Add(c.expire).Before(time.Now()) {
				if LOG {
					log.Println("Cache timeout Expired: " + k)
				}
				delete(c.items, k)
				i++
			}
		}
	}
	// 2. clean perc% of olders
	// How many do we have to clean ?
	total := len(c.items)
	num := total * perc / 100
	if LOG {
		log.Println("Quantity of elements to remove from cache:", num)
	}
	for i = 0; i < num; i++ {
		delete(c.items, c.pile[i])
	}
	c.pile = c.pile[i:]
	c.mutex.Unlock()
	return i
}

/* Verify:
  -------
  First, Clean(0) keeping all the entries, then deletes expired entries using Validator function.
  Returns the quantity of removed entries.
*/
func (c *XCache) Verify() int {
	// 1. clean all expired items, do not touch others
	i := c.Clean(0)
	// 2. If there is a validator, verifies anything
	if c.validator != nil {
		for k, x := range c.items {
			if b := c.validator(k, x.ctime); !b {
				if LOG {
					log.Println("Validator invalids entry: " + k)
				}
				c.mutex.Lock()
				delete(c.items, k)
				c.mutex.Unlock()
				i++
			}
		}
	}
	return i
}

/* Flush:
  ------
  Enpty the whole cache.
  Returns nothing.
*/
func (c *XCache) Flush() {
	c.mutex.Lock()
	// how to really deletes the data ? ( to free memory)
	c.items = make(map[string]*XCacheEntry)
	c.mutex.Unlock()
}

