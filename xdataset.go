package xcore

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// XDatasetDef is a special interface to implement a set of data that can be scanned recursively (by XTemplate for instance)
//   to search data into it, Stringify it, and set/get/del entries of data
//   The get* methods must accept a path id>id>id...
type XDatasetDef interface {
	// Stringify will dump the content into a human readable string
	fmt.Stringer   // please implement String()
	fmt.GoStringer // Please implement GoString()

	// Set will associate the data to the key. If it already exists, it will be replaced
	Set(key string, data interface{})

	// Get will return the value associated to the key if it exists, or bool = false
	Get(key string) (interface{}, bool)
	// Same as Get but will return the value associated to the key as a XDatasetDef if it exists, or bool = false
	GetDataset(key string) (XDatasetDef, bool)
	// Same as Get but will return the value associated to the key as a XDatasetCollectionDef if it exists, or bool = false
	GetCollection(key string) (XDatasetCollectionDef, bool)

	// Same as Get but will return the value associated to the key as a string if it exists, or bool = false
	GetString(key string) (string, bool)
	// Same as Get but will return the value associated to the key as a bool if it exists, or bool = false
	GetBool(key string) (bool, bool)
	// Same as Get but will return the value associated to the key as an int if it exists, or bool = false
	GetInt(key string) (int, bool)
	// Same as Get but will return the value associated to the key as a float64 if it exists, or bool = false
	GetFloat(key string) (float64, bool)
	// Same as Get but will return the value associated to the key as a Time if it exists, or bool = false
	GetTime(key string) (time.Time, bool)
	// Same as Get but will return the value associated to the key as a []String if it exists, or bool = false
	GetStringCollection(key string) ([]string, bool)
	// Same as Get but will return the value associated to the key as a []bool if it exists, or bool = false
	GetBoolCollection(key string) ([]bool, bool)
	// Same as Get but will return the value associated to the key as a []int if it exists, or bool = false
	GetIntCollection(key string) ([]int, bool)
	// Same as Get but will return the value associated to the key as a []float64 if it exists, or bool = false
	GetFloatCollection(key string) ([]float64, bool)
	// Same as Get but will return the value associated to the key as a []Time if it exists, or bool = false
	GetTimeCollection(key string) ([]time.Time, bool)

	// Del will delete the data associated to the key and deletes the key entry
	Del(key string)
	// Clone the object
	Clone() XDatasetDef
}

// =====================
// XDataset
// =====================

// XDataset is the basic interface dataset interface
type XDataset map[string]interface{}

// String will transform the XDataset into a readable string
func (d *XDataset) String() string {
	return d.GoString()
}

// GoString will transform the XDataset into a readable string for humans
func (d *XDataset) GoString() string {
	return fmt.Sprintf("%+v\n", *d)
}

// Set will add a variable key with value data to the XDataset
func (d *XDataset) Set(key string, data interface{}) {
	(*d)[key] = data
}

// Get will read the value of the key variable
func (d *XDataset) Get(key string) (interface{}, bool) {
	xid := strings.Split(key, ">")
	if len(xid) > 1 {
		subset, ok := (*d)[xid[0]]
		if !ok {
			return nil, false
		}
		if ds, ok := subset.(XDatasetDef); ok {
			return ds.Get(strings.Join(xid[1:], ">"))
		}
		if dsc, ok := subset.(XDatasetCollectionDef); ok {
			entry, err := strconv.Atoi(xid[1])
			if err != nil {
				return nil, false
			}
			ds, _ := dsc.Get(entry)
			if ds == nil {
				return nil, false
			}
			if len(xid) == 2 {
				return ds, true
			}
			return ds.Get(strings.Join(xid[2:], ">"))
		}
	}
	data, ok := (*d)[key]
	if ok {
		return data, true
	}
	return nil, false
}

// GetDataset will read the value of the key variable as a XDatasetDef cast type
func (d *XDataset) GetDataset(key string) (XDatasetDef, bool) {
	if val, ok := d.Get(key); ok {
		if val2, ok2 := val.(XDatasetDef); ok2 {
			return val2, true
		}
	}
	return nil, false
}

// GetCollection will read the value of the key variable as a XDatasetCollection cast type
func (d *XDataset) GetCollection(key string) (XDatasetCollectionDef, bool) {
	if val, ok := d.Get(key); ok {
		if val2, ok2 := val.(XDatasetCollectionDef); ok2 {
			return val2, true
		}
	}
	return nil, false
}

// GetString will read the value of the key variable as a string cast type
func (d *XDataset) GetString(key string) (string, bool) {
	if val, ok := d.Get(key); ok {
		return fmt.Sprint(val), true
	}
	return "", false
}

// GetBool will read the value of the key variable as a boolean cast type
// If the value is int, float, it will be convert with the rule 0: false, != 0: true
// If the value is anything else and it exists, it will return true if it's not nil
func (d *XDataset) GetBool(key string) (bool, bool) {
	if val, ok := d.Get(key); ok {
		if val2, ok2 := val.(bool); ok2 {
			return val2, true
		}
		if val2, ok2 := val.(int); ok2 {
			return val2 != 0, true
		}
		if val2, ok2 := val.(int8); ok2 {
			return val2 != 0, true
		}
		if val2, ok2 := val.(int16); ok2 {
			return val2 != 0, true
		}
		if val2, ok2 := val.(int32); ok2 {
			return val2 != 0, true
		}
		if val2, ok2 := val.(int64); ok2 {
			return val2 != 0, true
		}
		if val2, ok2 := val.(uint); ok2 {
			return val2 != 0, true
		}
		if val2, ok2 := val.(uint8); ok2 {
			return val2 != 0, true
		}
		if val2, ok2 := val.(uint16); ok2 {
			return val2 != 0, true
		}
		if val2, ok2 := val.(uint32); ok2 {
			return val2 != 0, true
		}
		if val2, ok2 := val.(uint64); ok2 {
			return val2 != 0, true
		}
		if val2, ok2 := val.(byte); ok2 {
			return val2 != 0, true
		}
		if val2, ok2 := val.(float32); ok2 {
			return val2 != 0, true
		}
		if val2, ok2 := val.(float64); ok2 {
			return val2 != 0, true
		}
		if val != nil {
			return true, true
		}
	}
	return false, false
}

// GetInt will read the value of the key variable as an integer cast type
// If the value is bool, will return 0/1
// If the value is float, will return integer part of value
func (d *XDataset) GetInt(key string) (int, bool) {
	if val, ok := d.Get(key); ok {
		if val2, ok2 := val.(bool); ok2 {
			if val2 {
				return 1, true
			}
			return 0, true
		}
		if val2, ok2 := val.(int); ok2 {
			return val2, true
		}
		if val2, ok2 := val.(int8); ok2 {
			return int(val2), true
		}
		if val2, ok2 := val.(int16); ok2 {
			return int(val2), true
		}
		if val2, ok2 := val.(int32); ok2 {
			return int(val2), true
		}
		if val2, ok2 := val.(int64); ok2 {
			return int(val2), true
		}
		if val2, ok2 := val.(uint); ok2 {
			return int(val2), true
		}
		if val2, ok2 := val.(uint8); ok2 {
			return int(val2), true
		}
		if val2, ok2 := val.(uint16); ok2 {
			return int(val2), true
		}
		if val2, ok2 := val.(uint32); ok2 {
			return int(val2), true
		}
		if val2, ok2 := val.(uint64); ok2 {
			return int(val2), true
		}
		if val2, ok2 := val.(byte); ok2 {
			return int(val2), true
		}
		if val2, ok2 := val.(float32); ok2 {
			return int(val2), true
		}
		if val2, ok2 := val.(float64); ok2 {
			return int(val2), true
		}
	}
	return 0, false
}

// GetFloat will read the value of the key variable as a float64 cast type
func (d *XDataset) GetFloat(key string) (float64, bool) {
	if val, ok := d.Get(key); ok {
		if val2, ok2 := val.(bool); ok2 {
			if val2 {
				return 1.0, true
			}
			return 0.0, true
		}
		if val2, ok2 := val.(int); ok2 {
			return float64(val2), true
		}
		if val2, ok2 := val.(int8); ok2 {
			return float64(val2), true
		}
		if val2, ok2 := val.(int16); ok2 {
			return float64(val2), true
		}
		if val2, ok2 := val.(int32); ok2 {
			return float64(val2), true
		}
		if val2, ok2 := val.(int64); ok2 {
			return float64(val2), true
		}
		if val2, ok2 := val.(uint); ok2 {
			return float64(val2), true
		}
		if val2, ok2 := val.(uint8); ok2 {
			return float64(val2), true
		}
		if val2, ok2 := val.(uint16); ok2 {
			return float64(val2), true
		}
		if val2, ok2 := val.(uint32); ok2 {
			return float64(val2), true
		}
		if val2, ok2 := val.(uint64); ok2 {
			return float64(val2), true
		}
		if val2, ok2 := val.(byte); ok2 {
			return float64(val2), true
		}
		if val2, ok2 := val.(float32); ok2 {
			return float64(val2), true
		}
		if val2, ok2 := val.(float64); ok2 {
			return val2, true
		}
	}
	return 0.0, false
}

// GetTime will read the value of the key variable as a time cast type
func (d *XDataset) GetTime(key string) (time.Time, bool) {
	if val, ok := d.Get(key); ok {
		if val2, ok2 := val.(time.Time); ok2 {
			return val2, true
		}
	}
	return time.Time{}, false
}

// GetStringCollection will read the value of the key variable as a collection of strings cast type
func (d *XDataset) GetStringCollection(key string) ([]string, bool) {
	if val, ok := d.Get(key); ok {
		if val2, ok2 := val.([]string); ok2 {
			return val2, true
		}
	}
	return nil, false
}

// GetBoolCollection will read the value of the key variable as a collection of bool cast type
func (d *XDataset) GetBoolCollection(key string) ([]bool, bool) {
	if val, ok := d.Get(key); ok {
		if val2, ok2 := val.([]bool); ok2 {
			return val2, true
		}
	}
	return nil, false
}

// GetIntCollection will read the value of the key variable as a collection of int cast type
func (d *XDataset) GetIntCollection(key string) ([]int, bool) {
	if val, ok := d.Get(key); ok {
		if val2, ok2 := val.([]int); ok2 {
			return val2, true
		}
	}
	return nil, false
}

// GetFloatCollection will read the value of the key variable as a collection of float cast type
func (d *XDataset) GetFloatCollection(key string) ([]float64, bool) {
	if val, ok := d.Get(key); ok {
		if val2, ok2 := val.([]float64); ok2 {
			return val2, true
		}
	}
	return nil, false
}

// GetTimeCollection will read the value of the key variable as a collection of time cast type
func (d *XDataset) GetTimeCollection(key string) ([]time.Time, bool) {
	if val, ok := d.Get(key); ok {
		if val2, ok2 := val.([]time.Time); ok2 {
			return val2, true
		}
	}
	return nil, false
}

// Del will deletes the variable
func (d *XDataset) Del(key string) {
	delete(*d, key)
}

// Clone will creates a totally new data memory cloned from this object
func (d *XDataset) Clone() XDatasetDef {
	cloned := &XDataset{}
	for id, val := range *d {
		clonedval := val
		// If the object is also cloneable, we clone it
		if cloneable1, ok := val.(interface{ Clone() XDatasetDef }); ok {
			clonedval = cloneable1.Clone()
		}
		if cloneable2, ok := val.(interface{ Clone() XDatasetCollectionDef }); ok {
			clonedval = cloneable2.Clone()
		}
		cloned.Set(id, clonedval)
	}
	return cloned
}

// XDatasetCollectionDef is the definition of a collection of XDatasetDefs
type XDatasetCollectionDef interface {
	// Stringify will dump the content into a human readable string
	fmt.Stringer   // please implement String()
	fmt.GoStringer // Please implement GoString()

	// Will add a datasetdef to the beginning of the collection
	Unshift(data XDatasetDef)
	// Will remove the first datasetdef of the collection and return it
	Shift() XDatasetDef

	// Will add a datasetdef to the end of the collection
	Push(data XDatasetDef)
	// Will remove the last datasetdef of the collection and return it
	Pop() XDatasetDef

	// Will count the quantity of entries
	Count() int

	// Will get the entry by the index and let it in the collection
	Get(index int) (XDatasetDef, bool)

	// Will search for the data associated to the key by priority (last entry is the most important)
	// returns bool = false if nothing has been found
	GetData(key string) (interface{}, bool)

	// Same as GetData but will convert the result to a string if possible
	// returns bool = false if nothing has been found
	GetDataString(key string) (string, bool)
	// Same as GetData but will convert the result to an int if possible
	// returns bool = false if nothing has been found
	GetDataInt(key string) (int, bool)
	// Same as GetData but will convert the result to a boolean if possible
	// returns second bool = false if nothing has been found
	GetDataBool(key string) (bool, bool)
	// Same as GetData but will convert the result to a float if possible
	// returns bool = false if nothing has been found
	GetDataFloat(key string) (float64, bool)
	// Same as GetData but will convert the result to a Time if possible
	// returns bool = false if nothing has been found
	GetDataTime(key string) (time.Time, bool)
	// Same as GetData but will convert the result to a XDatasetCollectionDef of data if possible
	// returns bool = false if nothing has been found
	GetCollection(key string) (XDatasetCollectionDef, bool)

	// Clone the object
	Clone() XDatasetCollectionDef
}

// =====================
// XDatasetConnection
// =====================

// XDatasetCollection is the basic collection of XDatasetDefs
type XDatasetCollection []XDatasetDef

// String will transform the XDataset into a readable string
func (d *XDatasetCollection) String() string {
	return d.GoString()
}

// GoString will transform the XDataset into a readable string for humans
func (d *XDatasetCollection) GoString() string {
	return fmt.Sprintf("%+v\n", *d)
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
