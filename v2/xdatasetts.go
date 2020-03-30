package xcore

import (
	"sync"
	"time"
)

// XDatasetTS is a thread safe xdataset (not thread safe) encapsulator
// XDatasetTS IS thread safe
type XDatasetTS struct {
	mutex sync.RWMutex
	data  XDatasetDef
}

// NewXDatasetTS builds a thread safe encapsulator on a XDataset compatible structure
func NewXDatasetTS(maindata XDatasetDef) *XDatasetTS {
	ds := &XDatasetTS{
		data: maindata,
	}
	return ds
}

// String will transform the XDataset into a readable string for humans
func (ds *XDatasetTS) String() string {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	return ds.data.String()
}

// GoString will transform the XDataset into a readable string for humans
func (ds *XDatasetTS) GoString() string {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	return ds.data.GoString()
}

// Set will add a variable key with value data to the XDatasetTS
func (ds *XDatasetTS) Set(key string, data interface{}) {
	ds.mutex.Lock()
	ds.data.Set(key, data)
	ds.mutex.Unlock()
}

// Get will read the value of the key variable
func (ds *XDatasetTS) Get(key string) (interface{}, bool) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	return ds.data.Get(key)
}

// GetDataset will read the value of the key variable as a XDatasetDef cast type
func (ds *XDatasetTS) GetDataset(key string) (XDatasetDef, bool) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	return ds.data.GetDataset(key)
}

// GetCollection will read the value of the key variable as a XDatasetCollection cast type
func (ds *XDatasetTS) GetCollection(key string) (XDatasetCollectionDef, bool) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	return ds.data.GetCollection(key)
}

// GetString will read the value of the key variable as a string cast type
func (ds *XDatasetTS) GetString(key string) (string, bool) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	return ds.data.GetString(key)
}

// GetBool will read the value of the key variable as a boolean cast type
// If the value is int, float, it will be convert with the rule 0: false, != 0: true
// If the value is anything else and it exists, it will return true if it's not nil
func (ds *XDatasetTS) GetBool(key string) (bool, bool) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	return ds.data.GetBool(key)
}

// GetInt will read the value of the key variable as an integer cast type
// If the value is bool, will return 0/1
// If the value is float, will return integer part of value
func (ds *XDatasetTS) GetInt(key string) (int, bool) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	return ds.data.GetInt(key)
}

// GetFloat will read the value of the key variable as a float64 cast type
func (ds *XDatasetTS) GetFloat(key string) (float64, bool) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	return ds.data.GetFloat(key)
}

// GetTime will read the value of the key variable as a time cast type
func (ds *XDatasetTS) GetTime(key string) (time.Time, bool) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	return ds.data.GetTime(key)
}

// GetStringCollection will read the value of the key variable as a collection of strings cast type
func (ds *XDatasetTS) GetStringCollection(key string) ([]string, bool) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	return ds.data.GetStringCollection(key)
}

// GetBoolCollection will read the value of the key variable as a collection of bool cast type
func (ds *XDatasetTS) GetBoolCollection(key string) ([]bool, bool) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	return ds.data.GetBoolCollection(key)
}

// GetIntCollection will read the value of the key variable as a collection of int cast type
func (ds *XDatasetTS) GetIntCollection(key string) ([]int, bool) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	return ds.data.GetIntCollection(key)
}

// GetFloatCollection will read the value of the key variable as a collection of float cast type
func (ds *XDatasetTS) GetFloatCollection(key string) ([]float64, bool) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	return ds.data.GetFloatCollection(key)
}

// GetTimeCollection will read the value of the key variable as a collection of time cast type
func (ds *XDatasetTS) GetTimeCollection(key string) ([]time.Time, bool) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	return ds.data.GetTimeCollection(key)
}

// Del will deletes the variable
func (ds *XDatasetTS) Del(key string) {
	ds.mutex.Lock()
	ds.data.Del(key)
	ds.mutex.Unlock()
}

// Clone will creates a totally new data memory cloned from this object
func (ds *XDatasetTS) Clone() XDatasetDef {
	cloned := &XDatasetTS{}
	ds.mutex.RLock()
	cloned.data = ds.data.Clone()
	ds.mutex.RUnlock()
	return cloned
}
