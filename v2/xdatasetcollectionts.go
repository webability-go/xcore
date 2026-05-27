package xcore

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// XDatasetCollectionTS is a thread-safe collection of XDatasetDefs.
// All methods are safe for concurrent use.
type XDatasetCollectionTS struct {
	mutex sync.RWMutex
	data  []XDatasetDef
}

// String will transform the XDatasetCollectionTS into a readable string
func (dc *XDatasetCollectionTS) String() string {
	str := "XDatasetCollectionTS["
	dc.mutex.RLock()
	for key, val := range dc.data {
		str += strconv.Itoa(key) + ":" + fmt.Sprint(val) + " "
	}
	dc.mutex.RUnlock()
	str += "]"
	return str
}

// GoString will transform the XDatasetCollectionTS into a readable string for humans
func (dc *XDatasetCollectionTS) GoString() string {
	return dc.String()
}

// Unshift will add a XDatasetDef at the beginning of the collection
func (dc *XDatasetCollectionTS) Unshift(data XDatasetDef) {
	dc.mutex.Lock()
	dc.data = append([]XDatasetDef{data}, dc.data...)
	dc.mutex.Unlock()
}

// Shift will remove and return the element at the beginning of the collection.
// Returns nil if the collection is empty.
func (dc *XDatasetCollectionTS) Shift() XDatasetDef {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()
	// 5.1/5.4: guard against empty-collection panic; defer ensures the lock is
	// always released even if a downstream operation panics.
	if len(dc.data) == 0 {
		return nil
	}
	data := dc.data[0]
	dc.data = dc.data[1:]
	return data
}

// Push will add a XDatasetDef at the end of the collection
func (dc *XDatasetCollectionTS) Push(data XDatasetDef) {
	dc.mutex.Lock()
	dc.data = append(dc.data, data)
	dc.mutex.Unlock()
}

// Pop will remove and return the element at the end of the collection.
// Returns nil if the collection is empty.
func (dc *XDatasetCollectionTS) Pop() XDatasetDef {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()
	// 5.1/5.4: guard against empty-collection panic; defer ensures the lock is
	// always released even if a downstream operation panics.
	if len(dc.data) == 0 {
		return nil
	}
	data := dc.data[len(dc.data)-1]
	dc.data = dc.data[:len(dc.data)-1]
	return data
}

// Count will return the quantity of elements into the collection
func (dc *XDatasetCollectionTS) Count() int {
	dc.mutex.RLock()
	defer dc.mutex.RUnlock()
	return len(dc.data)
}

// Get will retrieve an element by index from the collection
func (dc *XDatasetCollectionTS) Get(index int) (XDatasetDef, bool) {
	dc.mutex.RLock()
	defer dc.mutex.RUnlock()
	if index < 0 || index >= len(dc.data) {
		return nil, false
	}
	return dc.data[index], true
}

// GetData will retrieve the first available data identified by key from the collection
// searching from the last entry toward the first (last entry has highest priority).
//
// 5.2: a single RLock is held only while copying the slice header into a local
// snapshot; the lock is released before any element.Get calls.  This eliminates
// the TOCTOU race between the former separate length-read and element-access
// lock cycles, and avoids holding the collection lock while individual dataset
// Get methods run (which may themselves acquire locks in XDatasetTS entries).
func (dc *XDatasetCollectionTS) GetData(key string) (interface{}, bool) {
	dc.mutex.RLock()
	snapshot := make([]XDatasetDef, len(dc.data))
	copy(snapshot, dc.data)
	dc.mutex.RUnlock()
	for i := len(snapshot) - 1; i >= 0; i-- {
		val, ok := snapshot[i].Get(key)
		if ok {
			return val, true
		}
	}
	return nil, false
}

// GetDataString will retrieve the first available data identified by key and return it as a string.
//
// 5.5: nil values are consistently returned as ("", true) rather than being
// stringified to the literal text "nil", matching the behaviour of XDatasetCollection.GetDataString.
func (dc *XDatasetCollectionTS) GetDataString(key string) (string, bool) {
	v, ok := dc.GetData(key)
	if ok {
		if v == nil {
			return "", true
		}
		return fmt.Sprint(v), true
	}
	return "", false
}

// GetDataBool will retrieve the first available data identified by key and return it as a boolean
func (dc *XDatasetCollectionTS) GetDataBool(key string) (bool, bool) {
	if val, ok := dc.GetData(key); ok {
		if val2, ok2 := val.(bool); ok2 {
			return val2, true
		}
	}
	return false, false
}

// GetDataInt will retrieve the first available data identified by key and return it as an integer
func (dc *XDatasetCollectionTS) GetDataInt(key string) (int, bool) {
	if val, ok := dc.GetData(key); ok {
		if val2, ok2 := val.(int); ok2 {
			return val2, true
		}
	}
	return 0, false
}

// GetDataFloat will retrieve the first available data identified by key and return it as a float
func (dc *XDatasetCollectionTS) GetDataFloat(key string) (float64, bool) {
	if val, ok := dc.GetData(key); ok {
		if val2, ok2 := val.(float64); ok2 {
			return val2, true
		}
	}
	return 0, false
}

// GetDataTime will retrieve the first available data identified by key and return it as a time
func (dc *XDatasetCollectionTS) GetDataTime(key string) (time.Time, bool) {
	if val, ok := dc.GetData(key); ok {
		if val2, ok2 := val.(time.Time); ok2 {
			return val2, true
		}
	}
	return time.Time{}, false
}

// GetCollection will retrieve a collection identified by key from the XDatasetCollectionTS.
//
// 4.4: uses the two-value form of the type assertion so that a value stored
// under key that is not an XDatasetCollectionDef returns (nil, false) instead
// of panicking.
func (dc *XDatasetCollectionTS) GetCollection(key string) (XDatasetCollectionDef, bool) {
	v, ok := dc.GetData(key)
	if ok {
		if dsd, ok2 := v.(XDatasetCollectionDef); ok2 {
			return dsd, true
		}
	}
	return nil, false
}

// Clone will make a full copy of the object into memory.
//
// 5.3: holds RLock for the entire iteration so that concurrent Push/Pop/Shift/Unshift
// cannot mutate dc.data while we are copying element references.
func (dc *XDatasetCollectionTS) Clone() XDatasetCollectionDef {
	dc.mutex.RLock()
	defer dc.mutex.RUnlock()
	cloned := &XDatasetCollectionTS{}
	for _, val := range dc.data {
		cloned.data = append(cloned.data, val.Clone())
	}
	return cloned
}
