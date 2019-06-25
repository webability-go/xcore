package xcore

import (
  "fmt"
  "time"
//  "reflect"
)

/* 
The XDataset is a special interface to implement a set of data that can be scanned recursively (by XTemplate)
to search data into it, Stringify it, and set/get/del entries of data
*/

type XDatasetDef interface {
  // Stringify will dump the content into a human readable string
  Stringify() string

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

/* Basic dataset */
type XDataset map[string]interface{}

// makes an interface of XDataset to reuse for otrhe libraries and be sure we can call the functions
func (d *XDataset)Stringify() string {
  return fmt.Sprint(*d)
}

func (d *XDataset)Set(key string, data interface{}) {
  (*d)[key] = data
}

func (d *XDataset)Get(key string) (interface{}, bool) {
  data, ok := (*d)[key]
  if ok { return data, true }
  return nil, false
}

func (d *XDataset)GetDataset(key string) (XDatasetDef, bool) {
  if val, ok := (*d)[key]; ok {
    switch val.(type) {
      case XDatasetDef: return val.(XDatasetDef), true
    }
  }
  return nil, false
}

func (d *XDataset)GetCollection(key string) (XDatasetCollectionDef, bool) {
  if val, ok := (*d)[key]; ok {
    switch val.(type) {
      case XDatasetCollectionDef: return val.(XDatasetCollectionDef), true
    }
  }
  return nil, false
}

func (d *XDataset)GetString(key string) (string, bool) {
  data, ok := (*d)[key]
  if ok { return fmt.Sprint(data), true }
  return "", false
}

func (d *XDataset)GetBool(key string) (bool, bool) {
  if val, ok := (*d)[key]; ok {
    switch val.(type) {
      case bool: return val.(bool), true
      case int: return val.(int)!=0, true
      case float64: return val.(float64)!=0, true
    }
  }
  return false, false
}

func (d *XDataset)GetInt(key string) (int, bool) {
  if val, ok := (*d)[key]; ok {
    switch val.(type) {
      case int: return val.(int), true
      case float64: return int(val.(float64)), true
      case bool: if val.(bool) {
        return 1, true
      } else {
        return 0, true
      }
    }
  }
  return 0, false
}

func (d *XDataset)GetFloat(key string) (float64, bool) {
  if val, ok := (*d)[key]; ok {
    switch val.(type) {
      case float64: return val.(float64), true
      case int: return float64(val.(int)), true
      case bool: if val.(bool) {
        return 1.0, true
      } else {
        return 0.0, true
      }
    }
  }
  return 0, false
}

func (d *XDataset)GetTime(key string) (time.Time, bool) {
  if val, ok := (*d)[key]; ok {
    switch val.(type) {
      case time.Time: return val.(time.Time), true
    }
  }
  return time.Time{}, false
}

func (d *XDataset)GetStringCollection(key string) ([]string, bool) {
  if val, ok := (*d)[key]; ok {
    switch val.(type) {
      case []string: return val.([]string), true
    }
  }
  return nil, false
}

func (d *XDataset)GetBoolCollection(key string) ([]bool, bool) {
  if val, ok := (*d)[key]; ok {
    switch val.(type) {
      case []bool: return val.([]bool), true
    }
  }
  return nil, false
}

func (d *XDataset)GetIntCollection(key string) ([]int, bool) {
  if val, ok := (*d)[key]; ok {
    switch val.(type) {
      case []int: return val.([]int), true
    }
  }
  return nil, false
}

func (d *XDataset)GetFloatCollection(key string) ([]float64, bool) {
  if val, ok := (*d)[key]; ok {
    switch val.(type) {
      case []float64: return val.([]float64), true
    }
  }
  return nil, false
}

func (d *XDataset)GetTimeCollection(key string) ([]time.Time, bool) {
  if val, ok := (*d)[key]; ok {
    switch val.(type) {
      case []time.Time: return val.([]time.Time), true
    }
  }
  return nil, false
}

func (d *XDataset)Del(key string) {
  delete(*d, key)
}

// Check if we deep-clone the object ?
func (d *XDataset)Clone() XDatasetDef {
  cloned := &XDataset{}
  for id, val := range *d {
    clonedval := val
    // If the object is also cloneable, we clone it
    if cloneable, ok := val.(interface{Clone() XDatasetDef }); ok {
      clonedval = cloneable.Clone()
    }
    cloned.Set(id, clonedval)
  }
  return cloned
}

type XDatasetCollectionDef interface {
  // Stringify will dump the content into a human readable string
  Stringify() string

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

type XDatasetCollection []XDatasetDef

func (d *XDatasetCollection)Stringify() string {
  return fmt.Sprint(d)
}

func (d *XDatasetCollection)Unshift(data XDatasetDef) {
  *d = append([]XDatasetDef{data}, (*d)...)
}

func (d *XDatasetCollection)Shift() XDatasetDef {
  data := (*d)[0]
  *d = (*d)[1:]
  return data
}

func (d *XDatasetCollection)Push(data XDatasetDef) {
  *d = append(*d, data)
}

func (d *XDatasetCollection)Pop() XDatasetDef {
  data := (*d)[len(*d)-1]
  *d = (*d)[:len(*d)-1]
  return data
}

func (d *XDatasetCollection)Count() int {
  return len(*d)
}

func (d *XDatasetCollection)Get(index int) (XDatasetDef, bool) {
  if index < 0 || index >= len(*d) { return nil, false }
  return (*d)[index], true
}

func (d *XDatasetCollection)GetData(key string) (interface{}, bool) {
  for i := len(*d)-1; i >= 0; i-- {
    val, ok := (*d)[i].Get(key)
    if ok { return val, true }
  }
  return nil, false
}

func (d *XDatasetCollection)GetDataString(key string) (string, bool) {
  v, ok := d.GetData(key)
  if ok { return fmt.Sprint(v), true }
  return "", false
}

func (d *XDatasetCollection)GetDataBool(key string) (bool, bool) {
  if val, ok := d.GetData(key); ok {
    switch val.(type) {
      case bool: return val.(bool), true
    }
  }
  return false, false
}

func (d *XDatasetCollection)GetDataInt(key string) (int, bool) {
  if val, ok := d.GetData(key); ok {
    switch val.(type) {
      case int: return val.(int), true
    }
  }
  return 0, false
}

func (d *XDatasetCollection)GetDataFloat(key string) (float64, bool) {
  if val, ok := d.GetData(key); ok {
    switch val.(type) {
      case float64: return val.(float64), true
    }
  }
  return 0, false
}

func (d *XDatasetCollection)GetDataTime(key string) (time.Time, bool) {
  if val, ok := d.GetData(key); ok {
    switch val.(type) {
      case time.Time: return val.(time.Time), true
    }
  }
  return time.Time{}, false
}

func (d *XDatasetCollection)GetCollection(key string) (XDatasetCollectionDef, bool) {
  v, ok := d.GetData(key)
  // Verify v IS actually a XDatasetCollectionDef to avoid the error
  if ok { return v.(XDatasetCollectionDef), true }
  return nil, false
}

func (d *XDatasetCollection)Clone() XDatasetCollectionDef {
  cloned := &XDatasetCollection{}
  for _, val := range *d {
    *cloned = append(*cloned, val.Clone())
  }
  return cloned
}

