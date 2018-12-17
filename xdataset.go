package xcore

import (
  "fmt"
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
  // Same as Get but will return the value associated to the key as a string if it exists, or bool = false
  GetString(key string) (string, bool)
  // Same as Get but will return the value associated to the key as a XDatasetCollectionDef if it exists, or bool = false
  GetCollection(key string) (XDatasetCollectionDef, bool)
  // Del will delete the data associated to the key and deletes the key entry
  Del(key string)
}

type XDatasetCollectionDef interface {
  // Stringify will dump the content into a human readable string
  Stringify() string
  // Will add a datasetdef to the end of the collection
  Push(data XDatasetDef)
  // Will remove the last datasetdef of the collection and return it
  Pop() XDatasetDef
  // Will get the entry by the index and let it in the collection
  Get(index int) (XDatasetDef, bool)
  // Will count the quantity of entries
  Count() int
  // Will search for the data associated to the key by priority (last entry is the most important)
  // returns bool = false if nothing has been found
  GetData(key string) (interface{}, bool)
  // Same as GetData but will convert the result to a string
  // returns bool = false if nothing has been found
  GetDataString(key string) (string, bool)
  // Same as GetData but will convert the result to a collection of data
  // returns bool = false if nothing has been found
  GetDataRange(key string) (XDatasetCollectionDef, bool)
}

/* Basic data set for XTemplate */
type XDataset map[string]interface{}
type XDatasetCollection []XDatasetDef

// =====================
// XDataset
// =====================

// makes an interface of XDataset to reuse for otrhe libraries and be sure we can call the functions
func (d *XDataset)Stringify() string {
  return fmt.Sprint(*d)
}

func (d *XDataset)Get(key string) (interface{}, bool) {
  data, ok := (*d)[key]
  if ok { return data, true }
  return nil, false
}

func (d *XDataset)GetString(key string) (string, bool) {
  data, ok := (*d)[key]
  if ok { return fmt.Sprint(data), true }
  return "", false
}

func (d *XDataset)GetCollection(key string) (XDatasetCollectionDef, bool) {
  data, ok := (*d)[key]
  if ok { return data.(XDatasetCollectionDef), true }
  return nil, false
}

func (d *XDataset)Set(key string, data interface{}) {
  (*d)[key] = data
}

func (d *XDataset)Del(key string) {
  delete(*d, key)
}

// =====================
// XDatasetConnection
// =====================

func (d *XDatasetCollection)Stringify() string {
  return fmt.Sprint(d)
}

func (d *XDatasetCollection)Push(data XDatasetDef) {
  *d = append(*d, data)
}

func (d *XDatasetCollection)Pop() XDatasetDef {
  data := (*d)[len(*d)-1]
  *d = (*d)[:len(*d)-1]
  return data
}

func (d *XDatasetCollection)Get(index int) (XDatasetDef, bool) {
  if index < 0 || index >= len(*d) { return nil, false }
  return (*d)[index], true
}

func (d *XDatasetCollection)Count() int {
  return len(*d)
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

func (d *XDatasetCollection)GetDataRange(key string) (XDatasetCollectionDef, bool) {
  v, ok := d.GetData(key)
  // Verify v IS actually a XDatasetCollectionDef to avoid the error
  if ok { return v.(XDatasetCollectionDef), true }
  return nil, false
}

