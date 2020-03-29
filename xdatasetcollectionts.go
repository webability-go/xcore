package xcore

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// XDatasetCollectionTS is the basic collection of XDatasetDefs
type XDatasetCollectionTS struct {
	mutex sync.RWMutex
	data  []XDatasetDef
}

// String will transform the XDataset into a readable string
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

// GoString will transform the XDataset into a readable string for humans
func (dc *XDatasetCollectionTS) GoString() string {
	return dc.String()
}

// Unshift will adds a XDatasetDef at the beginning of the collection
func (dc *XDatasetCollectionTS) Unshift(data XDatasetDef) {
	dc.mutex.Lock()
	dc.data = append([]XDatasetDef{data}, dc.data...)
	dc.mutex.Unlock()
}

// Shift will remove the element at the beginning of the collection
func (dc *XDatasetCollectionTS) Shift() XDatasetDef {
	dc.mutex.Lock()
	data := dc.data[0]
	dc.data = dc.data[1:]
	dc.mutex.Unlock()
	return data
}

// Push will adds a XDatasetDef at the end of the collection
func (dc *XDatasetCollectionTS) Push(data XDatasetDef) {
	dc.mutex.Lock()
	dc.data = append(dc.data, data)
	dc.mutex.Unlock()
}

// Pop will remove the element at the end of the collection
func (dc *XDatasetCollectionTS) Pop() XDatasetDef {
	dc.mutex.Lock()
	data := dc.data[len(dc.data)-1]
	dc.data = dc.data[:len(dc.data)-1]
	dc.mutex.Unlock()
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

// GetData will retrieve the first available data identified by key from the collection ordered by index
func (dc *XDatasetCollectionTS) GetData(key string) (interface{}, bool) {
	dc.mutex.RLock()
	l := len(dc.data) - 1
	dc.mutex.RUnlock()
	for i := l; i >= 0; i-- {
		dc.mutex.RLock()
		dcc := dc.data[i]
		dc.mutex.RUnlock()
		val, ok := dcc.Get(key)
		if ok {
			return val, true
		}
	}
	return nil, false
}

// GetDataString will retrieve the first available data identified by key from the collection ordered by index and return it as a string
func (dc *XDatasetCollectionTS) GetDataString(key string) (string, bool) {
	v, ok := dc.GetData(key)
	if ok {
		return fmt.Sprint(v), true
	}
	return "", false
}

// GetDataBool will retrieve the first available data identified by key from the collection ordered by index and return it as a boolean
func (dc *XDatasetCollectionTS) GetDataBool(key string) (bool, bool) {
	if val, ok := dc.GetData(key); ok {
		if val2, ok2 := val.(bool); ok2 {
			return val2, true
		}
	}
	return false, false
}

// GetDataInt will retrieve the first available data identified by key from the collection ordered by index and return it as an integer
func (dc *XDatasetCollectionTS) GetDataInt(key string) (int, bool) {
	if val, ok := dc.GetData(key); ok {
		if val2, ok2 := val.(int); ok2 {
			return val2, true
		}
	}
	return 0, false
}

// GetDataFloat will retrieve the first available data identified by key from the collection ordered by index and return it as a float
func (dc *XDatasetCollectionTS) GetDataFloat(key string) (float64, bool) {
	if val, ok := dc.GetData(key); ok {
		if val2, ok2 := val.(float64); ok2 {
			return val2, true
		}
	}
	return 0, false
}

// GetDataTime will retrieve the first available data identified by key from the collection ordered by index and return it as a time
func (dc *XDatasetCollectionTS) GetDataTime(key string) (time.Time, bool) {
	if val, ok := dc.GetData(key); ok {
		if val2, ok2 := val.(time.Time); ok2 {
			return val2, true
		}
	}
	return time.Time{}, false
}

// GetCollection will retrieve a collection from the XDatasetCollectionTS
func (dc *XDatasetCollectionTS) GetCollection(key string) (XDatasetCollectionDef, bool) {
	v, ok := dc.GetData(key)
	// Verify v IS actually a XDatasetCollectionDef to avoid the error
	if ok {
		return v.(XDatasetCollectionDef), true
	}
	return nil, false
}

// Clone will make a full copy of the object into memory
func (dc *XDatasetCollectionTS) Clone() XDatasetCollectionDef {
	cloned := &XDatasetCollectionTS{}
	for _, val := range dc.data {
		cloned.data = append(cloned.data, val.Clone())
	}
	return cloned
}
