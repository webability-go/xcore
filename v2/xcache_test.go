package xcore

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// =============================================================================
// Examples (godoc)
// =============================================================================

func ExampleNewXCache() {
	cache1 := NewXCache("cacheid1", 0, 0)
	cache2 := NewXCache("cacheid2", 100, 0)
	cache3 := NewXCache("cacheid3", 0, 100*time.Minute)

	fmt.Println("All caches empty:", cache1.Count(), cache2.Count(), cache3.Count())
	// Output: All caches empty: 0 0 0
}

func ExampleXCache_Set() {
	cache := NewXCache("example", 0, 0)
	cache.Set("name", "WebAbility")
	v, invalid := cache.Get("name")
	fmt.Println(v, invalid)
	// Output: WebAbility false
}

func ExampleXCache_Get() {
	cache := NewXCache("example", 0, 0)
	// Get on a non existent entry returns (nil, false): the second value is the
	// "was invalidated" flag, not the "exists" flag.
	v, invalid := cache.Get("missing")
	fmt.Println(v, invalid)
	// Output: <nil> false
}

// =============================================================================
// Construction
// =============================================================================

func TestNewXCache(t *testing.T) {
	cache1 := NewXCache("cacheid1", 0, 0)
	cache2 := NewXCache("cacheid2", 100, 0)
	cache3 := NewXCache("cacheid3", 0, 100*time.Minute)

	if cache1 == nil || cache1.Count() != 0 {
		t.Fatal("Error creating NewXCache cache1")
	}
	if cache2 == nil || cache2.Count() != 0 {
		t.Fatal("Error creating NewXCache cache2")
	}
	if cache3 == nil || cache3.Count() != 0 {
		t.Fatal("Error creating NewXCache cache3")
	}

	// Metadata must be stored as given.
	if cache2.ID != "cacheid2" {
		t.Errorf("expected ID cacheid2, got %q", cache2.ID)
	}
	if cache2.Maxitems != 100 {
		t.Errorf("expected Maxitems 100, got %d", cache2.Maxitems)
	}
	if cache3.Expire != 100*time.Minute {
		t.Errorf("expected Expire 100m, got %v", cache3.Expire)
	}
	// A fresh cache must have a non-nil Validator-free state.
	if cache1.Validator != nil {
		t.Error("a fresh cache must have a nil Validator")
	}
}

// =============================================================================
// Set / Get / Del / Count
// =============================================================================

func TestXCacheSetGetDel(t *testing.T) {
	cache := NewXCache("cacheid", 0, 0)
	cache.Set("id1", "Data for id1")
	cache.Set("id2", "Data for id2")
	cache.Set("id3", "Data for id3")
	if cache.Count() != 3 {
		t.Fatalf("expected 3 entries, got %d", cache.Count())
	}

	data1, invalid1 := cache.Get("id1")
	if invalid1 || data1 != "Data for id1" {
		t.Fatalf("getting id1 failed: data=%v invalid=%v", data1, invalid1)
	}

	cache.Del("id1")
	if cache.Count() != 2 {
		t.Fatalf("expected 2 entries after delete, got %d", cache.Count())
	}

	// Get on a deleted/non existent key: (nil, false).
	data2, invalid2 := cache.Get("id1")
	if invalid2 || data2 != nil {
		t.Fatalf("getting deleted id1 should be (nil,false), got (%v,%v)", data2, invalid2)
	}
}

func TestXCacheGetNonExistent(t *testing.T) {
	cache := NewXCache("cacheid", 0, 0)
	v, invalid := cache.Get("nope")
	if v != nil || invalid {
		t.Fatalf("expected (nil,false) for missing key, got (%v,%v)", v, invalid)
	}
}

func TestXCacheDelNonExistent(t *testing.T) {
	cache := NewXCache("cacheid", 0, 0)
	// Must not panic and must not change the count.
	cache.Del("ghost")
	if cache.Count() != 0 {
		t.Fatalf("expected empty cache, got %d", cache.Count())
	}
}

func TestXCacheSetOverwrite(t *testing.T) {
	cache := NewXCache("cacheid", 0, 0)
	cache.Set("k", "v1")
	cache.Set("k", "v2")
	cache.Set("k", "v3")

	if cache.Count() != 1 {
		t.Fatalf("overwriting the same key must keep a single entry, got %d", cache.Count())
	}
	v, invalid := cache.Get("k")
	if invalid || v != "v3" {
		t.Fatalf("expected last value v3, got (%v,%v)", v, invalid)
	}
}

func TestXCacheStoresAnyType(t *testing.T) {
	cache := NewXCache("types", 0, 0)
	type rec struct {
		A int
		B string
	}
	cache.Set("int", 42)
	cache.Set("struct", rec{A: 1, B: "x"})
	cache.Set("slice", []int{1, 2, 3})
	cache.Set("nil", nil)

	if v, _ := cache.Get("int"); v != 42 {
		t.Errorf("int entry corrupted: %v", v)
	}
	if v, _ := cache.Get("struct"); v.(rec).B != "x" {
		t.Errorf("struct entry corrupted: %v", v)
	}
	// A nil value still counts as an existing, valid entry.
	v, invalid := cache.Get("nil")
	if invalid || v != nil {
		t.Errorf("nil value entry should be (nil,false), got (%v,%v)", v, invalid)
	}
	if cache.Count() != 4 {
		t.Errorf("expected 4 entries, got %d", cache.Count())
	}
}

// =============================================================================
// TTL (per entry)  -  SetTTL
// =============================================================================

func TestXCacheTTL(t *testing.T) {
	cache := NewXCache("ttl", 0, 0)
	cache.Set("k", "v")
	cache.SetTTL("k", 60*time.Millisecond)

	// Still valid right away.
	if v, invalid := cache.Get("k"); invalid || v != "v" {
		t.Fatalf("entry should be valid before TTL, got (%v,%v)", v, invalid)
	}

	time.Sleep(150 * time.Millisecond)

	// Expired by TTL: (nil, true) and the entry is removed.
	v, invalid := cache.Get("k")
	if !invalid || v != nil {
		t.Fatalf("entry expired by TTL should be (nil,true), got (%v,%v)", v, invalid)
	}
	if cache.Count() != 0 {
		t.Fatalf("expired TTL entry must be removed, count=%d", cache.Count())
	}
}

func TestXCacheSetTTLOnMissingDoesNothing(t *testing.T) {
	cache := NewXCache("ttl", 0, 0)
	// Must not panic and must not create an entry.
	cache.SetTTL("ghost", time.Second)
	if cache.Count() != 0 {
		t.Fatalf("SetTTL on a missing key must not create it, count=%d", cache.Count())
	}
}

// =============================================================================
// Global expiration (Expire)
// =============================================================================

func TestXCacheExpire(t *testing.T) {
	cache := NewXCache("cacheid", 0, 80*time.Millisecond)
	cache.Set("id1", "Data for id1")
	cache.Set("id2", "Data for id2")
	cache.Set("id3", "Data for id3")
	if cache.Count() != 3 {
		t.Fatalf("expected 3 entries, got %d", cache.Count())
	}

	time.Sleep(200 * time.Millisecond)

	// Count() does not actively purge: expiration is lazy, applied on Get().
	if cache.Count() != 3 {
		t.Fatalf("Count must not actively purge expired entries, got %d", cache.Count())
	}

	// Get an expired entry: (nil, true) and it gets removed.
	data1, invalid1 := cache.Get("id1")
	if !invalid1 || data1 != nil {
		t.Fatalf("expired id1 should be (nil,true), got (%v,%v)", data1, invalid1)
	}
	if cache.Count() != 2 {
		t.Fatalf("expected 2 entries after lazy purge, got %d", cache.Count())
	}

	data2, invalid2 := cache.Get("id2")
	if !invalid2 || data2 != nil {
		t.Fatalf("expired id2 should be (nil,true), got (%v,%v)", data2, invalid2)
	}
	if cache.Count() != 1 {
		t.Fatalf("expected 1 entry after lazy purge, got %d", cache.Count())
	}
}

// =============================================================================
// Validator hook
// =============================================================================

func TestXCacheValidatorValid(t *testing.T) {
	cache := NewXCache("val", 0, 0)
	cache.Validator = func(key string, ctime time.Time) bool { return true }
	cache.Set("k", "v")

	v, invalid := cache.Get("k")
	if invalid || v != "v" {
		t.Fatalf("validator returning true must keep the entry, got (%v,%v)", v, invalid)
	}
	if cache.Count() != 1 {
		t.Fatalf("entry must survive a valid check, count=%d", cache.Count())
	}
}

func TestXCacheValidatorInvalid(t *testing.T) {
	cache := NewXCache("val", 0, 0)
	cache.Validator = func(key string, ctime time.Time) bool { return false }
	cache.Set("k", "v")

	v, invalid := cache.Get("k")
	if !invalid || v != nil {
		t.Fatalf("validator returning false must invalidate: want (nil,true), got (%v,%v)", v, invalid)
	}
	if cache.Count() != 0 {
		t.Fatalf("invalidated entry must be removed, count=%d", cache.Count())
	}
}

// The validator hook must receive the entry key and its real creation time.
func TestXCacheValidatorReceivesKeyAndCtime(t *testing.T) {
	cache := NewXCache("val", 0, 0)

	var gotKey string
	var gotCtime time.Time
	cache.Validator = func(key string, ctime time.Time) bool {
		gotKey = key
		gotCtime = ctime
		return true
	}

	before := time.Now()
	cache.Set("the-key", "v")
	after := time.Now()
	cache.Get("the-key")

	if gotKey != "the-key" {
		t.Errorf("validator got key %q, expected the-key", gotKey)
	}
	if gotCtime.Before(before) || gotCtime.After(after) {
		t.Errorf("validator ctime %v outside the creation window [%v,%v]", gotCtime, before, after)
	}
}

// =============================================================================
// Clean / Flush
// =============================================================================

func TestXCacheClean(t *testing.T) {
	cache := NewXCache("cacheid", 110, 0)
	for i := 0; i < 100; i++ {
		cache.Set(fmt.Sprintf("id%d", i), i)
	}
	if cache.Count() != 100 {
		t.Fatalf("expected 100 entries, got %d", cache.Count())
	}

	// Clean 30% of the oldest -> removes 30, returns 30.
	removed := cache.Clean(30)
	if removed != 30 {
		t.Errorf("Clean(30) should report 30 removed, got %d", removed)
	}
	if cache.Count() != 70 {
		t.Fatalf("expected 70 entries after Clean(30), got %d", cache.Count())
	}

	// Clean 100% removes everything.
	removed = cache.Clean(100)
	if removed != 70 {
		t.Errorf("Clean(100) should report 70 removed, got %d", removed)
	}
	if cache.Count() != 0 {
		t.Fatalf("expected empty cache after Clean(100), got %d", cache.Count())
	}
}

func TestXCacheCleanZero(t *testing.T) {
	cache := NewXCache("cacheid", 0, 0)
	for i := 0; i < 10; i++ {
		cache.Set(fmt.Sprintf("id%d", i), i)
	}
	// Clean(0) on a cache with no expiration removes nothing.
	removed := cache.Clean(0)
	if removed != 0 {
		t.Errorf("Clean(0) should remove 0, got %d", removed)
	}
	if cache.Count() != 10 {
		t.Fatalf("Clean(0) must keep all entries, got %d", cache.Count())
	}
}

func TestXCacheCleanEmpty(t *testing.T) {
	cache := NewXCache("cacheid", 0, 0)
	// Cleaning an empty cache must not panic and must report 0.
	if removed := cache.Clean(50); removed != 0 {
		t.Errorf("cleaning an empty cache should remove 0, got %d", removed)
	}
}

// Clean must report the entries it purged by expiration, even with perc=0.
// (Regression: the previous implementation reset its counter and returned 0.)
func TestXCacheCleanReportsExpired(t *testing.T) {
	cache := NewXCache("clean-exp", 0, 60*time.Millisecond)
	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)
	time.Sleep(150 * time.Millisecond)

	removed := cache.Clean(0)
	if removed != 3 {
		t.Fatalf("Clean(0) should report 3 expired entries removed, got %d", removed)
	}
	if cache.Count() != 0 {
		t.Fatalf("expected empty cache after purging expired, got %d", cache.Count())
	}
}

// Expired entries count toward the perc% target: with 2 expired out of 10 and a
// 50% target (=5), only 3 more (the oldest) are removed -> 5 removed, 5 kept.
// This also guards pile/items consistency (a stale pile would leave Count != 5).
func TestXCacheCleanExpiredPlusPercent(t *testing.T) {
	cache := NewXCache("mix", 0, 60*time.Millisecond)
	cache.Set("old1", 1)
	cache.Set("old2", 2)
	time.Sleep(150 * time.Millisecond) // old1/old2 expire
	for i := 0; i < 8; i++ {
		cache.Set(fmt.Sprintf("new%d", i), i) // 8 fresh -> 10 total on entry
	}

	removed := cache.Clean(50)
	if removed != 5 {
		t.Fatalf("expected 5 removed (target 50%% of 10), got %d", removed)
	}
	if cache.Count() != 5 {
		t.Fatalf("expected 5 surviving entries, got %d", cache.Count())
	}
}

// When expiration alone meets or exceeds the perc% target, no extra (valid)
// entry must be evicted. Expired entries are still all purged (they are invalid).
func TestXCacheCleanExpiredExceedsTarget(t *testing.T) {
	cache := NewXCache("exceed", 0, 60*time.Millisecond)
	for i := 0; i < 5; i++ {
		cache.Set(fmt.Sprintf("e%d", i), i)
	}
	time.Sleep(150 * time.Millisecond) // all 5 expire

	// target = 20% of 5 = 1, but 5 are expired -> purge all 5, remove nothing extra.
	removed := cache.Clean(20)
	if removed != 5 {
		t.Fatalf("expected 5 removed (all expired), got %d", removed)
	}
	if cache.Count() != 0 {
		t.Fatalf("expected empty cache, got %d", cache.Count())
	}
}

// When expiration removes fewer than the target, the difference is taken from
// the oldest entries and the rest of the valid entries survive.
func TestXCacheCleanExpiredUnderTargetKeepsValid(t *testing.T) {
	cache := NewXCache("under", 0, 60*time.Millisecond)
	cache.Set("old", 1)
	time.Sleep(150 * time.Millisecond) // 1 expires
	for i := 0; i < 9; i++ {
		cache.Set(fmt.Sprintf("n%d", i), i) // 9 fresh -> 10 total on entry
	}

	// target = 40% of 10 = 4. expired = 1 -> 3 more removed -> 4 total, 6 kept.
	removed := cache.Clean(40)
	if removed != 4 {
		t.Fatalf("expected 4 removed, got %d", removed)
	}
	if cache.Count() != 6 {
		t.Fatalf("expected 6 surviving entries, got %d", cache.Count())
	}
}

func TestXCacheFlush(t *testing.T) {
	cache := NewXCache("cacheid", 0, 0)
	for i := 0; i < 50; i++ {
		cache.Set(fmt.Sprintf("id%d", i), i)
	}
	cache.Flush()
	if cache.Count() != 0 {
		t.Fatalf("expected empty cache after Flush, got %d", cache.Count())
	}
	// Cache must remain usable after a flush.
	cache.Set("again", "ok")
	if v, invalid := cache.Get("again"); invalid || v != "ok" {
		t.Fatalf("cache unusable after flush, got (%v,%v)", v, invalid)
	}
}

// Regression test for the security/consistency fix: Flush must also reset the
// internal pile. If the pile kept stale keys, a later Clean would slice the
// wrong keys and leave the cache in an inconsistent state (Count != 0 here).
func TestXCacheFlushResetsPile(t *testing.T) {
	cache := NewXCache("pile", 0, 0)
	for i := 0; i < 10; i++ {
		cache.Set(fmt.Sprintf("a%d", i), i)
	}
	cache.Flush()

	// Refill with fewer, brand new keys.
	cache.Set("b0", 0)
	cache.Set("b1", 1)
	cache.Set("b2", 2)

	// Clean(100) must wipe exactly the 3 current keys. With a stale pile this
	// would delete non existent keys and leave the 3 real entries behind.
	cache.Clean(100)
	if cache.Count() != 0 {
		t.Fatalf("pile not reset on Flush: Clean(100) left %d entries", cache.Count())
	}
}

// =============================================================================
// Maxitems auto-cleaning
// =============================================================================

func TestXCacheMaxitems(t *testing.T) {
	cache := NewXCache("max", 10, 0)
	for i := 0; i < 10; i++ {
		cache.Set(fmt.Sprintf("id%d", i), i)
	}
	// Exactly at the limit: no auto-clean yet.
	if cache.Count() != 10 {
		t.Fatalf("expected 10 entries at the limit, got %d", cache.Count())
	}

	// Exceeding the limit triggers an automatic Clean(10%) -> removes 1.
	cache.Set("id10", 10)
	if cache.Count() != 10 {
		t.Fatalf("exceeding Maxitems should auto-clean back to 10, got %d", cache.Count())
	}

	// The oldest entry (id0) must be the one evicted; the newest must remain.
	if _, ok := cache.Get("id10"); !ok {
		// Get returns the data; for a present valid entry "invalid" is false and
		// data is non-nil. We just check it is still retrievable.
	}
	if v, invalid := cache.Get("id10"); invalid || v != 10 {
		t.Errorf("newest entry id10 should remain valid, got (%v,%v)", v, invalid)
	}
}

// =============================================================================
// Verify (Validator-based scan)
// =============================================================================

func TestXCacheVerifyWithValidator(t *testing.T) {
	cache := NewXCache("ver", 0, 0)
	cache.Set("keep1", "1")
	cache.Set("drop", "2")
	cache.Set("keep2", "3")

	// Invalidate only the "drop" entry.
	cache.Validator = func(key string, ctime time.Time) bool { return key != "drop" }

	removed := cache.Verify()
	if removed != 1 {
		t.Errorf("Verify should report 1 invalidated entry, got %d", removed)
	}
	if cache.Count() != 2 {
		t.Fatalf("expected 2 surviving entries, got %d", cache.Count())
	}
	if v, invalid := cache.Get("keep1"); invalid || v != "1" {
		t.Errorf("keep1 should survive Verify, got (%v,%v)", v, invalid)
	}
	if v, invalid := cache.Get("drop"); invalid || v != nil {
		t.Errorf("drop should be gone (nil,false), got (%v,%v)", v, invalid)
	}
}

func TestXCacheVerifyNoValidator(t *testing.T) {
	cache := NewXCache("ver", 0, 60*time.Millisecond)
	cache.Set("a", 1)
	cache.Set("b", 2)
	time.Sleep(150 * time.Millisecond)

	// With no Validator, Verify still purges expired entries and reports them.
	removed := cache.Verify()
	if removed != 2 {
		t.Fatalf("Verify should report 2 expired entries removed, got %d", removed)
	}
	if cache.Count() != 0 {
		t.Fatalf("Verify must purge expired entries, got %d", cache.Count())
	}
}

// =============================================================================
// Concurrency / thread-safety  (run with: go test -race)
//
// Exercises the locking fixes: Get() updates rtime under the write lock, and
// Verify() iterates over a snapshot instead of the live map. Under -race any
// data race or "concurrent map read/write" would fail this test.
// =============================================================================

func TestXCacheConcurrentAccess(t *testing.T) {
	cache := NewXCache("conc", 0, 0)
	cache.Validator = func(key string, ctime time.Time) bool { return true }

	const workers = 40
	const ops = 300

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < ops; i++ {
				key := fmt.Sprintf("k%d", (w+i)%50)
				switch i % 6 {
				case 0:
					cache.Set(key, i)
				case 1:
					cache.Get(key)
				case 2:
					cache.Del(key)
				case 3:
					cache.Count()
				case 4:
					if i%30 == 0 {
						cache.Clean(10)
					}
				case 5:
					if i%50 == 0 {
						cache.Verify()
					}
				}
			}
		}(w)
	}
	wg.Wait()

	// The cache must remain consistent and usable afterwards.
	cache.Flush()
	if cache.Count() != 0 {
		t.Fatalf("cache not empty after Flush, got %d", cache.Count())
	}
	cache.Set("final", "ok")
	if v, invalid := cache.Get("final"); invalid || v != "ok" {
		t.Fatalf("cache unusable after concurrent stress, got (%v,%v)", v, invalid)
	}
}

// Concurrent Get on the same key stresses the rtime/pile update path that now
// runs under the write lock.
func TestXCacheConcurrentGetSameKey(t *testing.T) {
	cache := NewXCache("conc2", 0, 0)
	cache.Set("hot", "value")

	const workers = 50
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				if v, invalid := cache.Get("hot"); invalid || v != "value" {
					t.Errorf("unexpected read: (%v,%v)", v, invalid)
					return
				}
			}
		}()
	}
	wg.Wait()

	if cache.Count() != 1 {
		t.Fatalf("expected the single hot entry to survive, got %d", cache.Count())
	}
}
