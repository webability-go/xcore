/*
	The XCore package is a set of basic objects for programation (XCache, XDataset, XLanguage, XTemplate)
	For GO, the actual existing code includes:

		- XCache: Application Memory Caches for any purpose,
		- XDataset: Basic nested data structures for any purpose (template injection, configuration files, database records, etc),
		- XLanguage: language dependant text tables,
		- XTemplate: template system with meta language.

	The Package hast been used for years on professional PHP projects in the WebAbility Core for PHP program.
*/
// XCache is a library to cache all the data you want into current application memory for a very fast access to the data.
// The access to the data support multithreading and concurrency. For the same reason, this type of cache is not persistant (if you exit the application)
// and cannot grow too much (as memory is the limit).
// However, you can control a timeout of each cache piece, and eventually the comparison against a source (file, database, etc) to invalid the cache.
// 
// -----------------------
// 1. Overview
// 
// Declare a new XCache with NewXCache()
// 
// Then you can use the 3 basic functions to control the content of the cache: Get/Set/Del.
// You can put any kind of data into your XCache.
// The XCache is thread safe.
// 
// The cache can be limited in quantity of entries and timeout for data. The cache is automanaged (for invalid expired data) and can be cleaned partially or totally manually.
// 
// If you want some stats of the cache, use the Count function.
// 
// Example:
// 
// import "github.com/webability-go/xcore"
// 
// var myfiles = xcore.NewXCache("myfiles", 0, 0)
// 
// func usemycache() {
//   myfiles.Set("https://developers.webability.info:82/", "somedata")
//   myfiles.Set("/home/sites/file2.txt", "someotherdata")
// 
//   go somefunc()
//   
//   fmt.Println("Quantity of data into cache:", myfiles.count())
// }
// 
// func somefunc() {
//   data, invalid := myfiles.Get("https://developers.webability.info:82/");
//   
// }

package xcore

// VERSION: is the used version nombre of the XCore library.
const VERSION = "0.1.1"

// if LOG is set to TRUE, LOG indicates to the XCore libraries to log a trace of functions called, with most important parameters.
//
// LOG can be set to true or false dynamically to trace only parts of code on demand.
var LOG = false
