package xcore

import (
	"log"
	"sync"
	"time"
)

// XCacheEntry is the cache basic structure to save some data in memory.
type XCacheEntry struct {
	// The cache entry has a time to measure expiration if needed, or time of entry in cache:
	// - ctime is the creation time (used to validate the object against its source).
	ctime time.Time
	// - rtime is the last read time (used to clean the cache: the less accessed objects are removed).
	rtime time.Time
	// - ttl is the max duration of object in cache (ctime + ttl > now = invalid)
	ttl time.Duration
	// The data as itself is an interface to whatever the user need to cache.
	data interface{}
}

// XCache is the main cache structure, that contains a collection of XCacheEntries and some metadata.
type XCache struct {
	// "ID": XCache has a unique id (informative).
	ID string
	// "Maxitems": The user can creates a cache with a maximum number of elements into it. In this case, when the cache reaches the maximum number of elements stored, then the system makes a clean of 10% of the oldest elements. This type of use is not recommended since is it heavy in CPU use to clean the cache.
	Maxitems int
	// "Validator" is a function that can be set to check the validity of the data (for instance if the data originates from a file or a database). The validator is called for each Get (and can be heavy for CPU or can wait a long time, for instance if the check is an external database on another cluster). Beware of this.
	Validator func(string, time.Time) bool
	// "Expire": The user can also create an expiration duration, so every elements in the cache is invalidated after a certain amount of time. It is more recommended to use the cache with an expiration duration. The obsolete objects are destroyed when the user tries to use them and return a "non existence" on Get. (this does not use CPU or extra locks).
	Expire time.Duration
	// Not available from outside for security, access of data is based on a mutex
	// "mutex": The cache owns a mutex to lock access to data to read/write/delete/clean the data, to allow concurrency and multithreading of the cache.
	mutex sync.RWMutex
	// "pile": The pile keeps the "ordered by date of reading" object keys, so it's fast to clean the data.
	items map[string]*XCacheEntry
	// "items": The items are a map to cache entries, acceved by the key of entries.
	pile []string
}

// NewXCache function will create a new XCache structure.
// The XCache is resident in memory, supports multithreading and concurrency.
// "id" is the unique id of the XCache.
// "maxitems" is the max authorized quantity of objects into the XCache. If 0, the cache hast no limit in quantity of objects.
// "expire" is a max duration of the objects into the cache. If 0, no limit
// Returns the *XCache created.
func NewXCache(id string, maxitems int, expire time.Duration) *XCache {
	if LOG {
		log.Printf("Creating cache with data {id: %s, maxitems: %d, expire: %d}", id, maxitems, expire)
	}
	return &XCache{
		ID:        id,
		Maxitems:  maxitems,
		Validator: nil,
		Expire:    expire,
		items:     make(map[string]*XCacheEntry),
	}
}

// Set will set an entry in the cache.
// If the entry already exists, just replace it with a new creation date.
// If the entry does not exist, it will insert it in the cache and if the cache if full (maxitems reached), then a clean is called to remove 10%.
// Returns nothing.
func (c *XCache) Set(key string, indata interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// check if the entry already exists
	_, ok := c.items[key]
	c.items[key] = &XCacheEntry{ctime: time.Now(), rtime: time.Now(), data: indata}
	if ok {
		c.removeFromPile(key)
	}
	c.pile = append(c.pile, key)
	if c.Maxitems > 0 && len(c.items) >= c.Maxitems {
		// We need a cleaning
		c.mClean(10)
	}
}

// Set will set a TTL on the entry in the cache.
// If the entry exists, just ads the TTL to the entry
// If the entry does not exist, it does nothing
// Returns nothing.
func (c *XCache) SetTTL(key string, duration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, ok := c.items[key]
	if ok {
		c.items[key].ttl = duration
	}
}

// removeFromPile will remove an entry key from the ordered pile.
// No lock into this func since it has been set by entry func already
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

// mGet will get the entry of the cache.
func (c *XCache) mGet(key string) (*XCacheEntry, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	x, ok := c.items[key]
	return x, ok
}

// Get will get the value of an entry.
// If the entry does not exists, returns nil, false.
// If the entry exists and is invalidated by time or validator function, then returns nil, true.
// If the entry is good, return <value>, false.
func (c *XCache) Get(key string) (interface{}, bool) {
	if x, ok := c.mGet(key); ok {
		// expired by TTL ?
		if x.ttl != 0 {
			if x.ctime.Add(x.ttl).Before(time.Now()) {
				c.Del(key)
				/*
					c.mutex.Lock()
					delete(c.items, key)
					c.removeFromPile(key)
					c.mutex.Unlock()
				*/
				return nil, true
			}
		}
		if c.Validator != nil {
			if b := c.Validator(key, x.ctime); !b {
				if LOG {
					log.Println("Validator invalids entry: " + key)
				}
				c.Del(key)
				/*
					c.mutex.Lock()
					delete(c.items, key)
					c.removeFromPile(key)
					c.mutex.Unlock()
				*/
				return nil, true
			}
		}
		// expired ?
		if c.Expire != 0 {
			if x.ctime.Add(c.Expire).Before(time.Now()) {
				if LOG {
					log.Println("Cache timeout Expired: " + key)
				}
				c.Del(key)
				/*
					c.mutex.Lock()
					delete(c.items, key)
					c.removeFromPile(key)
					c.mutex.Unlock()
				*/
				return nil, true
			}
		}
		x.rtime = time.Now()
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.removeFromPile(key)
		c.pile = append(c.pile, key)
		return x.data, false
	}
	return nil, false
}

// Del will delete the entry of the cache if it exists.
func (c *XCache) Del(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.items, key)
	// we should check if the entry exists before trying to removing
	c.removeFromPile(key)
}

// Count will return the quantity of entries in the cache.
func (c *XCache) Count() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	x := len(c.items)
	return x
}

// mClean will delete expired entries, and free perc% of max items based on time. It does not set locks. The caller must do it
// perc = 0 to 100 (percentage to clean).
// Returns quantity of removed entries.
// It Will **not** verify the cache against its source (if Validator is set). If you want to scan that, use the Verify function.
func (c *XCache) mClean(perc int) int {
	if LOG {
		log.Println("Cleaning cache")
	}
	i := 0
	// 1. clean all expired items
	if c.Expire != 0 {
		for k, x := range c.items {
			if x.ctime.Add(c.Expire).Before(time.Now()) {
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
	return i
}

// Clean will delete expired entries, and free perc% of max items based on time.
// perc = 0 to 100 (percentage to clean).
// Returns quantity of removed entries.
// It Will **not** verify the cache against its source (if Validator is set). If you want to scan that, use the Verify function.
func (c *XCache) Clean(perc int) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.mClean(perc)
}

// Verify will first, Clean(0) keeping all the entries, then will delete expired entries using the Validator function.
// Returns the quantity of removed entries.
// Based on what the validator function does, calling Verify can be **very** slow and cpu dependant. Be very careful.
func (c *XCache) Verify() int {
	// 1. clean all expired items, do not touch others
	i := c.Clean(0)
	// 2. If there is a validator, verifies anything
	if c.Validator != nil {
		for k, x := range c.items {
			if b := c.Validator(k, x.ctime); !b {
				if LOG {
					log.Println("Validator invalids entry: " + k)
				}
				c.Del(k)
				/*
					c.mutex.Lock()
					delete(c.items, k)
					c.removeFromPile(k)
					c.mutex.Unlock()
				*/
				i++
			}
		}
	}
	return i
}

// Flush will empty the whole cache and free all the memory of it.
// Returns nothing.
func (c *XCache) Flush() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// how to really deletes the data ? ( to free memory)
	c.items = make(map[string]*XCacheEntry)
}
