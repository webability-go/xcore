package xcore

import (
  "fmt"
)

/* 
The XDataset is a special interface to implement a set of data that can be scanned recursively (by XTemplate)
to search data into it, print it, and set/get/del entries of data
*/

type XDatasetDef interface {
  // Print will dump the content into a human readable string
  Print() string
  // Set will associate the data to the key. If it already exists, it will be replaced
  Set(key string, data interface{})
  // Get will return the value associated to the key if it exists, or nil
  Get(key string) interface{}
  // Same as Get but will return the value associated to the key as a string if it exists, or ""
  GetString(key string) string
  // Same as Get but will return the value associated to the key as a XDatasetCollectionDef if it exists, or nil
  GetCollection(key string) XDatasetCollectionDef
  // Del will delete the data associated to the key and deletes the key entry
  Del(key string)
}

type XDatasetCollectionDef interface {
  // Print will dump the content into a human readable string
  Print() string
  // Will add a datasetdef to the end of the collection
  Push(data XDatasetDef)
  // Will remove the last datasetdef of the collection and return it
  Pop() XDatasetDef
  // Will get the entry by the index and let it in the collection
  Get(index int) XDatasetDef
  // Will count the quantity of entries
  Count() int
  // Will search for the data associated to the key by priority (last entry is the most important)
  // returns nil if nothing has been found
  GetData(key string) interface{}
  // Same as GetData but will convert the result to a string
  // returns "" if nothing has been found
  GetDataString(key string) string
  GetDataRange(key string) XDatasetCollectionDef
}

/* Basic data set for XTemplate */
type XDataset map[string]interface{}
type XDatasetCollection []XDatasetDef

// =====================
// XDataset
// =====================

// makes an interface of XDataset to reuse for otrhe libraries and be sure we can call the functions
func (d XDataset)Print() string {
  return fmt.Sprint(d)
//  return fmt.Sprintf("DIRECCION DEL OBJETO: %p %p ", &d, d)
}

func (d XDataset)Get(key string) interface{} {
  data, ok := (d)[key]
  if ok { return data }
  return nil
}

func (d XDataset)GetString(key string) string {
  data, ok := (d)[key]
  if ok { return fmt.Sprint(data) }
  return ""
}

func (d XDataset)GetCollection(key string) XDatasetCollectionDef {
  data, ok := (d)[key]
  if ok { return data.(XDatasetCollectionDef) }
  return nil
}

func (d XDataset)Set(key string, data interface{}) {
  (d)[key] = data
}

func (d XDataset)Del(key string) {
  delete(d, key)
}

// =====================
// XDatasetConnection
// =====================

func (d XDatasetCollection)Print() string {
  return fmt.Sprint(d)
//  return fmt.Sprintf("DATOS DEL PRIMER ELEMENTO: %p ", d.Get(0))
}

func (d XDatasetCollection)Push(data XDatasetDef) {
  d = append(d, data)
}

func (d XDatasetCollection)Pop() XDatasetDef {
  data := d[len(d)-1]
  d = d[:len(d)-1]
  return data
}

func (d XDatasetCollection)Get(index int) XDatasetDef {
  return d[index]
}

func (d XDatasetCollection)Count() int {
  return len(d)
}

func (d XDatasetCollection)GetData(key string) interface{} {
  for i := len(d)-1; i >= 0; i-- {
    val := d[i].Get(key)
    if val != nil {
      return val
    }
  }
  return nil
}

func (d XDatasetCollection)GetDataString(key string) string {
  v := d.GetData(key)
  if v != nil {
//    return fmt.Sprintf("DATOS DEL PRIMER ELEMENTO: %p %p ", v, &v)
    return fmt.Sprint(v)
  }
  return ""
}

func (d XDatasetCollection)GetDataRange(key string) XDatasetCollectionDef {
  v := d.GetData(key)
  
  fmt.Println(v)
  return nil
  
  if v != nil { return v.(XDatasetCollectionDef) }
  return nil
}

