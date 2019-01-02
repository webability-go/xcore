@UTF-8

XCore for GO v0
=============================

The XCore package is used to build basic object for programation. for the WebAbility compatility code
For GO, the actual existing code includes:
- Memory Caches
- Templates
- Language tables

TO DO:
======
- XCache: comments in code, manual
- Apply XDataset for XConfig
- XCache: makes some test, what is faster, 10000x go threads sleeping one with each data into the thread and a channel to wake them up and communicate the data, or like it is right now (mutex and concurrent acceses for a memory dynamic map for 10000 memory pointers)
- XLanguage comments in code and manual
- XTemplate must concatenate strings after compilation
- Implements functions as data entry for template Execute (simple data or loop funcions, can get backs anything, creates an interface)
- Implements 2 parameters for &&, 3 parameters for @@ and ??
- Implements templates derivation (.first, .last, .#num, .keyvalue, .none, etc)
- XCache: activate persistant cache too (shared memory)

Version Changes Control
=======================

V0.0.3 - 2018-??-??
-----------------------
- Added XCache.Flush function
- Function XCache.Del implemented

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

1. Overview
-----------------------

Declare a new XCache with NewXCache()

Then you can use the 3 basic functions to control the content of the cache: Get/Set/Del.

If you want some stats of the cache, use the Count function.

Example:

import "github.com/webability-go/xcache"

var myfiles := xfile.NewXCache("myfiles", 0, true, 0)

func usemycache() {
  myfiles.Set("/home/sites/file1.txt", "somedata")
  myfiles.Set("/home/sites/file2.txt", "somedata")

  go somefunc()
}


2. Reference
-----------------------

To use the package: 

import "github.com/webability-go/xcache"

List of functions:

func NewXCache(id string, maxitems int, isfile bool, expire time.Duration) *XCache
------------------------
Creates a new XCache structure. 
The XCache is resident in memory, supports multithreading and concurrency.
"id" is the unique id of the XCache. 
maxitems is the max authorized quantity of objects into the XCache.
isfile is a boolean set to true then the entries IDs are filepath on hard disk and the system will check expiration date against the file last modif date. If the file is newer than the cache, the entry is invalidated and need recalculation.
expire is a max duration of the objects into the cache.


func (c *XCache)Set(key string, indata interface{}) 
------------------------
Set an entry into our cache. If there was already an entry with the same key, it will be replaced.

func (c *XCache)Get(key string) (interface{}, bool)
------------------------
Gets the value of the entry. 
- If the entry exists and is valid, returns the pointer to the object and false.
- If the entry exists and is not valid anymore, returns nil and true.
- If the entry does not exist, return nil and false.

func (c *XCache)Del(key string)
------------------------
Deletes the entry in the XCache

func (c *XCache)Clean(perc int) int
------------------------
Cleans the cache to be able to receive more data: will first scan to invalidate and deletes expired entries (if there is an expiration time), then free 10% of items limit (if any)
It Will **not** verify the cache against its source (if isfile is set to true). If you want to scan that, use the Verify function.

func (c *XCache)Verify() int
------------------------
Will scan the cache to first clean it (keeping 100% of the valid data), then scan all the data against its source if any.

func (c *XCache)Count(key string) int
------------------------
Gets the quantity of valid entries into the cache.

func (c *XCache)Flush()
------------------------
Invalidate the whole cache and empty it


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
