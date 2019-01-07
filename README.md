@UTF-8

XCore for GO v0
=============================

The XCore package is used to build basic object for programation. for the WebAbility compatility code
For GO, the actual existing code includes:
- XCache: Application Memory Caches
- XDataset: Basic nested data structures for any purpose (template injection, configuration files, database records, etc)
- XLanguage: language dependant text tables
- XTemplate: template systeme with meta language

TO DO:
======
- XLanguage comments in code and manual
- XTemplate must concatenate strings after compilation
- Implements functions as data entry for template Execute (simple data or loop funcions, can get backs anything, creates an interface)
- Implements 2 parameters for &&, 3 parameters for @@ and ??
- Implements templates derivation (.first, .last, .#num, .keyvalue, .none, etc)

- XCache: makes some test, what is faster, 10000x go threads sleeping one with each data into the thread and a channel to wake them up and communicate the data, or like it is right now (mutex and concurrent acceses for a memory dynamic map for 10000 memory pointers)
- XCache: activate persistant cache too (shared memory) ?????


Version Changes Control
=======================

V0.0.5 - 2019-??-??
-----------------------
- 

V0.0.4 - 2019-01-02
-----------------------
- XDataset.Get* functions added to comply with any type of data of a dataset for templates, config, database record etc.

V0.0.3 - 2019-01-02
-----------------------
- Added XCache.Flush function
- Function XCache.Del implemented
- Function XCache.Clean implemented for expiration, and free some space
- Function XCache.Verify created
- Function XCache.SetValidator added, to check cache validity agains a validator function
- Files flags and code removed from XCache. If the cache is a file, the user should controls the files with its own Validator function (original funcions put in examples as a file validator). This lets a lots of flexibility to validate against any source of data (files, database, complex calculations, external streams, etc)
- XCache is ready for candidate release

V0.0.2 - 2018-12-17
-----------------------
- Creation of XCache with all set of functions.
- Creation of XLanguage with all set of functions.
- Creation of XTemplate with all set of functions. Basic work done
- Creation of a set of interfaces that XTemplate need to execute and inject the template, 
- Creation of a basic XDataset and colection based on interfaces to build a set of data for the template.
- Added xcore.go with version number as constant

V0.0.1 - 2018-11-14
-----------------------
- First basic commit with XLanguage object created



Manual:
=======================

XCache
=======================
XCache is a library to cache all the data you want into current application memory for a very fast access to the data.
The access to the data support multithreading and concurrency. For the same reason, this type of cache is not persistant (if you exit the application)
and cannot grow too much (as memory is the limit).
However, you can control a timeout of each cache piece, and eventually the comparison against a source file to invalid the cache.

-----------------------
1. Overview

Declare a new XCache with NewXCache()

Then you can use the 3 basic functions to control the content of the cache: Get/Set/Del.

If you want some stats of the cache, use the Count function.

Example:

```
import "github.com/webability-go/xcache"

var myfiles := xfile.NewXCache("myfiles", 0, 0)

func usemycache() {
  myfiles.Set("/home/sites/file1.txt", "somedata")
  myfiles.Set("/home/sites/file2.txt", "somedata")

  go somefunc()
}
```


-----------------------
2. Reference

To use the package: 

import "github.com/webability-go/xcache"

List of types:

XCacheEntry:
------------------------
  The cache entry has a time to measure expiration if needed, or time of entry in cache.
  - ctime is creation time (used to validate the object against its source).
  - rtime is last read time (used to clean the cache: the less accessed objects are removed).
  The data as itself is an interface to whatever the user need to cache.


XCache:
------------------------
  The XCache has an id (informative).
  - The user can creates a cache with a maximum number of elements if need. In this case, when the cache reaches the maximum number of elements stored, then the system makes a clean of 10% of oldest elements. This type of use is not recommended since is it heavy in CPU use to clean the cache.
  - The user can also create an expiration duration, so every elements in the cache is invalidated after a certain amount of time. It is more recommended to use the cache with an expiration duration. The obsolete objects are destroyed when the user tries to use them and return a "non existance" on Get. (this does not use CPU or extra locks.
  - The Validator is a function that can be set to check the validity of the data (for instance if the data originates from a file or a database). The validator is called for each Get (and can be heavy for CPU or can wait a long time, for instance if the check is an external database on another cluster). Beware of this.
  - The cache owns a mutex to lock access to data to read/write/delete/clean the data, to allow concurrency and multithreading of the cache.
  - The pile keeps the "ordered by date of reading" object keys, so it's fast to clean the data.
  - Finally, the items are a map to cache entries, acceved by the key of entries.

  
List of functions:

func NewXCache(id string, maxitems int, expire time.Duration) *XCache
------------------------
  Creates a new XCache structure. 
  The XCache is resident in memory, supports multithreading and concurrency.
  "id" is the unique id of the XCache. 
  maxitems is the max authorized quantity of objects into the XCache.
  expire is a max duration of the objects into the cache.
  Returns the *XCache created.


func (c *XCache)SetValidator(f func(string, time.Time) bool)
------------------------
  Sets the validator function to check every entry in the cache against its original source, for each Get and Verify calls.
  Returns nothing.


func (c *XCache)Set(key string, indata interface{}) 
------------------------
  Sets an entry in the cache.
  If the entry already exists, just replace it with a new creation date.
  If the entry does not exist, it will insert it in the cache and if the cache if full (maxitems reached), then a clean is called to remove 10%.
  Returns nothing.


func (c *XCache)Get(key string) (interface{}, bool)
------------------------
  get the value of an entry.
  If the entry does not exists, returns nil, false
  If the entry exists and is invalidated by time or validator function, then returns nil, true
  If the entry is good, return <value>, false


func (c *XCache)Del(key string)
------------------------
  deletes the entry of the cache if it exists.


func (c *XCache)Count(key string) int
------------------------
  returns the quantity of entries in the cache.


func (c *XCache)Clean(perc int) int
------------------------
  deletes expired entries, and free perc% of max items based on time.
  perc = 0 to 100 (percentage to clean).
  Returns quantity of removed entries.
  It Will **not** verify the cache against its source (if isfile is set to true). If you want to scan that, use the Verify function.


func (c *XCache)Verify() int
------------------------
  First, Clean(0) keeping all the entries, then deletes expired entries using Validator function.
  Returns the quantity of removed entries.

func (c *XCache)Flush()
------------------------
  Empty the whole cache.
  Returns nothing.





XLanguage
=======================

1. Overview
------------------------

The XLanguage table of text entries can be loaded from XML file, XML string or normal file or string.

The XML Format is:
<?xml version="1.0" encoding="UTF-8"?>
<language id="NAMEOFLANGUAGE" lang="LG">
  <entry id="ENTRYNAME">ENTRYVALUE</entry>
  <entry id="ENTRYNAME">ENTRYVALUE</entry>
</language>
where NAMEOFLANGUAGE is the name of your table entry, for example "homepate", "user_report", etc
      LG is the ISO-3369 2 letters language ID, for example "es" for spanish, "en" for english
      ENTRYNAME is the ID of the entry, for example "greating", "yourname", "submitbutton"
      ENTRYVALUE is the text for your entry, for example "Hello", "You are:", "Save" if your table is in english

The Flat Text format is:
ENTRYNAME=ENTRYVALUE
ENTRYNAME=ENTRYVALUE
where ENTRYNAME is the ID of the entry, for example "greating", "yourname", "submitbutton"
      ENTRYVALUE is the text for your entry, for example "Hello", "You are:", "Save" if your table is in english

There is no name of table or language in this format (you "know" what you are loading)

The advantage to use XML forma is to have more control over your language, and eventyally add attributes into your entry, 
for instance translated="yes/no", verified="yes/no", and any other data that your system could insert

Create a new XLanguage empty structure:

- NewXLanguage

There are 4 functions to create the language from a file or string, flat text or XML text:

- NewXLanguageFromXMLFile
- NewXLanguageFromXMLString
- NewXLanguageFromFile
- NewXLanguageFromString

Then you can use the set of basic access functions:

- Set/Get/Del/SetName/SetLanguage/GetName/GetLanguage

2. Reference
------------------------






XTemplate
=======================

1. Overview
------------------------

class to compile and keep a Template string
A template is a set of HTML/XML (or any other language) set with a meta language made of:

Comments:
   %-- comments --%
Fields:
  {{field}}
  {{field>Subfield>Subfield}}
Language injection
  ##entry##
Subtemplates:
   xml/html code
   [[id]]
     xml/html code
     [[id]]
       xml/html code indented
     [[]]
     xml/html code
   [[]]
Meta elements:
   ??xx??   conditions
   @@xx@@   loops
   &&xx&&   references
   !!xx!!   debug (dump)

  2. Meta Language Reference
------------------------


  3. Functions Reference
------------------------




---
