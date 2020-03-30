package xcore

import (
	"fmt"
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
