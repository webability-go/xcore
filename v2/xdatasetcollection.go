package xcore

import (
	"fmt"
	"strconv"
	"time"
)

// =====================
// XDatasetCollection
// =====================

// XDatasetCollection is the basic collection of XDatasetDefs
type XDatasetCollection []XDatasetDef

// String will transform the XDataset into a readable string
func (d *XDatasetCollection) String() string {
	str := "XDatasetCollection["
	for key, val := range *d {
		str += strconv.Itoa(key) + ":" + fmt.Sprint(val) + " "
	}
	str += "]"
	return str
}

// GoString will transform the XDataset into a readable string for humans
func (d *XDatasetCollection) GoString() string {
	return d.String()
}

// Unshift will adds a XDatasetDef at the beginning of the collection
func (d *XDatasetCollection) Unshift(data XDatasetDef) {
	*d = append([]XDatasetDef{data}, (*d)...)
}

// Shift will remove the element at the beginning of the collection
func (d *XDatasetCollection) Shift() XDatasetDef {
	data := (*d)[0]
	*d = (*d)[1:]
	return data
}

// Push will adds a XDatasetDef at the end of the collection
func (d *XDatasetCollection) Push(data XDatasetDef) {
	*d = append(*d, data)
}

// Pop will remove the element at the end of the collection
func (d *XDatasetCollection) Pop() XDatasetDef {
	data := (*d)[len(*d)-1]
	*d = (*d)[:len(*d)-1]
	return data
}

// Count will return the quantity of elements into the collection
func (d *XDatasetCollection) Count() int {
	return len(*d)
}

// Get will retrieve an element by index from the collection
func (d *XDatasetCollection) Get(index int) (XDatasetDef, bool) {
	if index < 0 || index >= len(*d) {
		return nil, false
	}
	return (*d)[index], true
}

// GetData will retrieve the first available data identified by key from the collection ordered by index
func (d *XDatasetCollection) GetData(key string) (interface{}, bool) {
	for i := len(*d) - 1; i >= 0; i-- {
		val, ok := (*d)[i].Get(key)
		if ok {
			return val, true
		}
	}
	return nil, false
}

// GetDataString will retrieve the first available data identified by key from the collection ordered by index and return it as a string
func (d *XDatasetCollection) GetDataString(key string) (string, bool) {
	v, ok := d.GetData(key)
	if ok {
		return fmt.Sprint(v), true
	}
	return "", false
}

// GetDataBool will retrieve the first available data identified by key from the collection ordered by index and return it as a boolean
func (d *XDatasetCollection) GetDataBool(key string) (bool, bool) {
	if val, ok := d.GetData(key); ok {
		if val2, ok2 := val.(bool); ok2 {
			return val2, true
		}
	}
	return false, false
}

// GetDataInt will retrieve the first available data identified by key from the collection ordered by index and return it as an integer
func (d *XDatasetCollection) GetDataInt(key string) (int, bool) {
	if val, ok := d.GetData(key); ok {
		if val2, ok2 := val.(int); ok2 {
			return val2, true
		}
	}
	return 0, false
}

// GetDataFloat will retrieve the first available data identified by key from the collection ordered by index and return it as a float
func (d *XDatasetCollection) GetDataFloat(key string) (float64, bool) {
	if val, ok := d.GetData(key); ok {
		if val2, ok2 := val.(float64); ok2 {
			return val2, true
		}
	}
	return 0, false
}

// GetDataTime will retrieve the first available data identified by key from the collection ordered by index and return it as a time
func (d *XDatasetCollection) GetDataTime(key string) (time.Time, bool) {
	if val, ok := d.GetData(key); ok {
		if val2, ok2 := val.(time.Time); ok2 {
			return val2, true
		}
	}
	return time.Time{}, false
}

// GetCollection will retrieve a collection from the XDatasetCollection
func (d *XDatasetCollection) GetCollection(key string) (XDatasetCollectionDef, bool) {
	v, ok := d.GetData(key)
	// Verify v IS actually a XDatasetCollectionDef to avoid the error
	if ok {
		return v.(XDatasetCollectionDef), true
	}
	return nil, false
}

// Clone will make a full copy of the object into memory
func (d *XDatasetCollection) Clone() XDatasetCollectionDef {
	cloned := &XDatasetCollection{}
	for _, val := range *d {
		*cloned = append(*cloned, val.Clone())
	}
	return cloned
}
