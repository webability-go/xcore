package xcore

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// =====================
// XDataset
// =====================

// XDataset is the basic dataset type that is xdatasetdef inferface compatible
// XDataset IS NOT thread safe
type XDataset map[string]interface{}

// NewXDataset is used to build an XDataset from a standard map
func NewXDataset(data map[string]interface{}) XDatasetDef {
	// Scan data and encapsulate it into the XDataset
	ds := &XDataset{}
	for i, v := range data {
		switch v.(type) {
		case []interface{}:
			dsc := &XDatasetCollection{}
			for _, it := range v.([]interface{}) {
				dscit, ok := it.(map[string]interface{})
				if ok {
					dscc := NewXDataset(dscit)
					dsc.Push(dscc)
				}
			}
			ds.Set(i, dsc)
		case []map[string]interface{}:
			ds.Set(i, NewXDatasetCollection(v.([]map[string]interface{})))
		case map[string]interface{}:
			ds.Set(i, NewXDataset(v.(map[string]interface{})))
		default:
			ds.Set(i, v)
		}
	}
	return ds
}

// String will transform the XDataset into a readable string for humans
func (d *XDataset) String() string {
	sdata := []string{}
	for key, val := range *d {
		sdata = append(sdata, key+":"+fmt.Sprintf("%v", val))
	}
	sort.Strings(sdata) // Lets be sure the print is always the same presentation
	return "xcore.XDataset{" + strings.Join(sdata, " ") + "}"
}

// GoString will transform the XDataset into a readable string for humans
func (d *XDataset) GoString() string {
	sdata := []string{}
	for key, val := range *d {
		sdata = append(sdata, key+":"+fmt.Sprintf("%#v", val))
	}
	sort.Strings(sdata) // Lets be sure the print is always the same presentation
	return "#xcore.XDataset{" + strings.Join(sdata, " ") + "}"
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
		if val == nil {
			return "", true
		}
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
		if val2, ok2 := val.(time.Time); ok2 {
			return !val2.Equal(time.Time{}), true
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
		// byte is alias of uint8, no need for it
		//		if val2, ok2 := val.(byte); ok2 {
		//			return val2 != 0, true
		//		}
		if val2, ok2 := val.(float32); ok2 {
			return val2 != 0, true
		}
		if val2, ok2 := val.(float64); ok2 {
			return val2 != 0, true
		}
		return val != nil, true
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
		//		if val2, ok2 := val.(byte); ok2 {
		//			return int(val2), true
		//		}
		if val2, ok2 := val.(float32); ok2 {
			return int(val2), true
		}
		if val2, ok2 := val.(float64); ok2 {
			return int(val2), true
		}
		if val2, ok2 := val.(time.Time); ok2 {
			if val2.Equal(time.Time{}) {
				return 0, true
			}
			return int(val2.Unix()), true
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
		//		if val2, ok2 := val.(byte); ok2 {
		//			return float64(val2), true
		//		}
		if val2, ok2 := val.(float32); ok2 {
			return float64(val2), true
		}
		if val2, ok2 := val.(float64); ok2 {
			return val2, true
		}
		if val2, ok2 := val.(time.Time); ok2 {
			if val2.Equal(time.Time{}) {
				return 0, true
			}
			return float64(val2.Unix()), true
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
		// int, float, conversion ?
	}
	return time.Time{}, false
}

// GetStringCollection will read the value of the key variable as a collection of strings cast type
func (d *XDataset) GetStringCollection(key string) ([]string, bool) {
	if val, ok := d.Get(key); ok {
		if val2, ok2 := val.([]string); ok2 {
			return val2, true
		}
		// other types: conversion?
	}
	return nil, false
}

// GetBoolCollection will read the value of the key variable as a collection of bool cast type
func (d *XDataset) GetBoolCollection(key string) ([]bool, bool) {
	if val, ok := d.Get(key); ok {
		if val2, ok2 := val.([]bool); ok2 {
			return val2, true
		}
		// other types: conversion?
	}
	return nil, false
}

// GetIntCollection will read the value of the key variable as a collection of int cast type
func (d *XDataset) GetIntCollection(key string) ([]int, bool) {
	if val, ok := d.Get(key); ok {
		if val2, ok2 := val.([]int); ok2 {
			return val2, true
		}
		// other types: conversion?
	}
	return nil, false
}

// GetFloatCollection will read the value of the key variable as a collection of float cast type
// If the field is not an array it will be converted to an array.
// If the field is an array of another type, it will be converted.
func (d *XDataset) GetFloatCollection(key string) ([]float64, bool) {
	if val, ok := d.Get(key); ok {
		if val2, ok2 := val.([]float64); ok2 {
			return val2, true
		}
		// other types: conversion?

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
