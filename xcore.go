// Copyright Philippe Thomassigny 2004-2019
// Use of this source code is governed by a MIT licence
// license that can be found in the LICENSE file.

// Package xcore is a set of basic objects for programation (XCache, XDataset, XLanguage, XTemplate).
// For GO, the actual existing code includes:
//
// - XCache: Application Memory Caches for any purpose,
//
// - XDataset: Basic nested data structures for any purpose (template injection, configuration files, database records, etc),
//
// - XLanguage: language dependant text tables,
//
// - XTemplate: template system with meta language.
//
// The Package hast been used for years on professional PHP projects in the WebAbility Core for PHP program and is now available for GO.
// It is already used on sites that serve more than 60 million pages a month.
//
// XCache
//
// XCache is a library to cache all the data you want into current application memory for a very fast access to the data.
// The access to the data support multithreading and concurrency. For the same reason, this type of cache is not persistent (if you exit the application)
// and cannot grow too much (as memory is the limit).
// However, you can control a timeout of each cache piece, and eventually a comparison function against a source (file, database, etc) to invalid the cache.
//
// 1. Declare a new XCache with NewXCache() function:
//
//  import "github.com/webability-go/xcore"
//
//  // 50 items max
//  var myfiles = xcore.NewXCache("myfiles", 50, 0)
//
//  // 10 minutes expiration
//  var mydbtable = xcore.NewXCache("mydb-table", 0, 10 * time.Minute)
//
//  // No direct limits on the cache
//  var myotherdbtable = xcore.NewXCache("mydb-table", 0, 0)
//
// 2. Once you have declared the cache, you can fill it with anything you want. The main cache object is an interface{}
// so you can put here anything you need, from simple variables to complex structures. You need to use the Set funcion:
//
//  func main() {
//    myfiles.Set("https://developers.webability.info/", "somedata")
//    myfiles.Set("/home/sites/file2.txt", "someotherdata")
//  	myrecords := GetAllMyDatabaseTableData()
//  	for _, rec := range myrecords {
//      key, _ := rec.GetString("key")
//  		mydbtable.Set(key, rec)
//    }
//  }
//
// 3. To use the cache, just ask for your entry with Get function:
//
//  func usemycache() {
//
//    filedata, invalid := myfiles.Get("https://developers.webability.info/");
//    dbdata, invalid2 := mydbtable.Get("4455");
//    // do something with the fetched data
//  }
//
// 4. To maintain the cache you may need Del function, to delete a specific entry (maybe because you deleted the record in database).
// You may also need Clean function to deletes a percentage of the cache, or Flush to deletes it all.
// The Verify function is used to check cache entries against their sources through the Validator function.
// Be very carefull, if the cache is big or the Validator function is complex (maybe ask for a remote server information),
// the verification may be VERY slow and huge CPU use.
// The Count function gives some stats about the cache.
//
//  func analyze() {
//
//    // Actual size of cache
//    fmt.Println(mydbtable.Count())
//    // Removes 30% of objects
//    objectsremoved := mydbtable.Clean(30)
//    // New size of cache
//    fmt.Println(mydbtable.Count())
//    // totally flush the cache
//    mydbtable.Flush()
//    // New size of cache, should be 0
//    fmt.Println(mydbtable.Count())
//  }
//
// The XCache is thread safe.
// The cache can be limited in quantity of entries and timeout for data. The cache is automanaged (for invalid expired data) and can be cleaned partially or totally manually.
//
//
// XLanguage
//
// The XLanguage table of text entries can be loaded from XML file, XML string or normal file or string.
// XML Format is:
//  <?xml version="1.0" encoding="UTF-8"?>
//  <language id="NAMEOFLANGUAGE" lang="LG">
//    <entry id="ENTRYNAME">ENTRYVALUE</entry>
//    <entry id="ENTRYNAME">ENTRYVALUE</entry>
//  </language>
//
//  where NAMEOFLANGUAGE is the name of your table entry, for example "homepate", "user_report", etc
//       LG is the ISO-3369 2 letters language ID, for example "es" for spanish, "en" for english
//       ENTRYNAME is the ID of the entry, for example "greating", "yourname", "submitbutton"
//       ENTRYVALUE is the text for your entry, for example "Hello", "You are:", "Save" if your table is in english
//
// Flat Text format is:
//  ENTRYNAME=ENTRYVALUE
//  ENTRYNAME=ENTRYVALUE
//
// where ENTRYNAME is the ID of the entry, for example "greating", "yourname", "submitbutton"
//       ENTRYVALUE is the text for your entry, for example "Hello", "You are:", "Save" if your table is in english
//
// There is no name of table or language in this format (you "know" what you are loading)
//
// The advantage to use XML forma is to have more control over your language, and eventyally add attributes into your entry,
// for instance translated="yes/no", verified="yes/no", and any other data that your system could insert
//
// Create a new XLanguage empty structure:
//
// - NewXLanguage()
//
// There are 4 functions to create the language from a file or string, flat text or XML text:
//
// - NewXLanguageFromXMLFile()
//
// - NewXLanguageFromXMLString()
//
// - NewXLanguageFromFile()
//
// - NewXLanguageFromString()
//
// Then you can use the set of basic access functions:
//
// - Set/Get/Del/SetName/SetLanguage/GetName/GetLanguage
//
package xcore

// VERSION is the used version nombre of the XCore library.
const VERSION = "0.2.1"

// LOG is the flag to activate logging on the library.
// if LOG is set to TRUE, LOG indicates to the XCore libraries to log a trace of functions called, with most important parameters.
// LOG can be set to true or false dynamically to trace only parts of code on demand.
var LOG = false
