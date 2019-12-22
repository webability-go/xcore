@UTF-8

[![Go Report Card](https://goreportcard.com/badge/github.com/webability-go/xcore)](https://goreportcard.com/report/github.com/webability-go/xcore)
[![GoDoc](https://godoc.org/github.com/webability-go/xcore?status.png)](https://godoc.org/github.com/webability-go/xcore)
[![GolangCI](https://golangci.com/badges/github.com/webability-go/xcore.svg)](https://golangci.com)

XCore for GO v0
=============================

The XCore package is used to build basic object for programation. for the WebAbility compatility code
For GO, the actual existing code includes:
- XCache: Application Memory Caches
- XDataset: Basic nested data structures for any purpose (template injection, configuration files, database records, etc)
- XLanguage: language dependant text tables
- XTemplate: template systeme with meta language

Manuals are available on godoc.org [![GoDoc](https://godoc.org/github.com/webability-go/xcore?status.png)](https://godoc.org/github.com/webability-go/xcore)


TO DO:
======
- XLanguage comments in code and manual
- XTemplate comments in code and manual
- XTemplate must concatenate strings after compilation
- Implements functions as data entry for template Execute (simple data or loop funcions, can get backs anything, creates an interface)
- Implements 2 parameters for &&, 3 parameters for @@ and ??
- Implements templates derivation (.first, .last, .#num, .keyvalue, .none, etc)

- Some improvements to check, later:
- XCache: activate persistant cache too (shared memory) ????? maybe not for go itself, but for instance to talk with other memory data used by other languages and apps, or to not loose the caches if the app is restarted.


Version Changes Control
=======================

v0.2.2 - 2019-12-21
- XLanguage now support golang x/text/language instead of direct iso 2 charater language
- godoc manuals for xlanguage, xdataset and xtemplate prepared

v0.2.1 - 2019-12-13
- XCache manual enhanced with examples

v0.2.0 - 2019-12-06
- XCache Code simplified to expose XCache definition as public, remove not usefull funcion (Get*)
- XCache 0.2.0 is not compatible with XCache 0.1.* , you may need to change your code
- Added more conversions between int-float-bool in XDataset.Get*

v0.1.2 - 2019-12-05
- Code cleaned to meet golangci standards, golint checks, more documentation.

V0.1.1 - 2019-11-05
-----------------------
- XCore Code comments enhanced to publish in godoc.org as libraries documentation

V0.1.0 - 2019-10-18
-----------------------
- Code cleaned to pass 100% of goreportcard.com. Card note added in this document

V0.0.9 - 2019-07-18
-----------------------
- Error corrected on XCache: removing an element from a slice when the element is the last one was causing out of bound index.
- XCache.maxitem = 0 (no number of elements limit) is corrected: it was not working

V0.0.8 - 2019-06-25
-----------------------
- Added Clone on XDatasetDef and XDataCollectionsetDef
- XDataset testunit added

V0.0.7 - 2019-03-06
-----------------------
- Time functions added to XDatasetDef and XDatasetCollectionDef interfaces, and XDataset and XDatasetCollection structures
- Manual for XCache finished
- Manual for XDataset finished
- Preformat for XLanguage manual
- Preformat for XTemplate manual

V0.0.6 - 2019-02-07
-----------------------
- Added xcache.GetId(), xcache.GetMax() and xcache.GetExpire()
- XCache Documentation modified

V0.0.5 - 2019-02-05
-----------------------
- Added conversion between types con XDataset.Get* functions
- Manuals for XDataSet and XTemplate complemented

V0.0.4 - 2019-01-06
-----------------------
- XDataset.Get* functions added to comply with any type of data of a dataset for templates, config, database record etc.
- XCache manual completed.

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



func NewXCache(id string, maxitems int, expire time.Duration) *XCache
------------------------
  Creates a new XCache structure.
  The XCache is resident in memory, supports multithreading and concurrency.
  "id" is the unique id of the XCache.
  maxitems is the max authorized quantity of objects into the XCache.
  expire is a max duration of the objects into the cache.
  Returns the *XCache created.


func (c *XCache)GetId() string
------------------------
  exposes the ID of the cache


func (c *XCache)GetMax() int
------------------------
  exposes the max quantity of items of the cache


func (c *XCache)GetExpire() time.Duration
------------------------
  exposes the expiration time of the cache


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
  It Will **not** verify the cache against its source (if Validator is set). If you want to scan that, use the Verify function.


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

```
<?xml version="1.0" encoding="UTF-8"?>
<language id="NAMEOFLANGUAGE" lang="LG">
  <entry id="ENTRYNAME">ENTRYVALUE</entry>
  <entry id="ENTRYNAME">ENTRYVALUE</entry>
</language>
```

where:
- NAMEOFLANGUAGE is the name of your table entry, for example "homepage", "user_report", etc
- LG is the ISO-3369 2 letters language ID, for example "es" for spanish, "en" for english
- ENTRYNAME is the ID of the entry, for example "greating", "yourname", "submitbutton"
- ENTRYVALUE is the text for your entry, for example "Hello", "You are:", "Save" if your table is in english

The Flat Text format is:

```
ENTRYNAME=ENTRYVALUE
ENTRYNAME=ENTRYVALUE
```

where:
- ENTRYNAME is the ID of the entry, for example "greating", "yourname", "submitbutton"
- ENTRYVALUE is the text for your entry, for example "Hello", "You are:", "Save" if your table is in english

There is no name of table or language in this format (you "know" what you are loading)

The advantage to use XML forma is to have more control over your language, and eventyally add attributes into your entry, for instance translated="yes/no", verified="yes/no", and any other data that your system could insert

Create a new XLanguage empty structure:

- NewXLanguage

There are 4 functions to create the language from a file or string, flat text or XML text:

- NewXLanguageFromXMLFile
- NewXLanguageFromXMLString
- NewXLanguageFromFile
- NewXLanguageFromString

Then you can use the set of basic access functions:

- Set/Get/Del/SetName/SetLanguage/GetName/GetLanguage

Example:

```
lang, err := xcore.NewXLanguageFromXMLString(`
<?xml version="1.0" encoding="UTF-8"?>
<language id="language-demo" lang="en">
  <entry id="entry1">Welcome to</entry>
  <entry id="entry2">XCore</entry>
</language>
`)

tr1 = lang.Get("entry1")
tr2 = lang.Get("entry2")

fmt.Println(tr1, tr2)
```


2. Reference
------------------------

To use the package:

import "github.com/webability-go/xcore"

List of types:

XLanguage:
------------------------
  The language entry has a name and a language as main parameters.
  The Entries map structure contains all the different language entries  key=value


func NewXLanguage(name string, lang string) *XLanguage {
   Creates an empty Language structure with a name and a language

func NewXLanguageFromXMLFile(file string) (*XLanguage, error) {
   Creates an XLanguage structure with the data from the XML file
   Returns nil if there is an error

func NewXLanguageFromXMLString(xml string) (*XLanguage, error) {

func NewXLanguageFromFile(file string) (*XLanguage, error) {

func NewXLanguageFromString(data string) (*XLanguage, error) {

/* LoadXMLFile:
   Loads a language from an XML file and replace the content of the XLanguage structure with the new data
*/
func (l *XLanguage)LoadXMLFile(file string) error {

func (l *XLanguage)LoadXMLString(data string) error {

func (l *XLanguage)LoadFile(file string) error {

func (l *XLanguage)LoadString(data string) error {

func (l *XLanguage)SetName(name string) {

func (l *XLanguage)SetLanguage(lang string) {

func (l *XLanguage)GetName() string {

func (l *XLanguage)GetLanguage() string {

func (l *XLanguage)Set(entry string, value string) {

func (l *XLanguage)Get(entry string) string {

func (l *XLanguage)Del(entry string) {

func (l *XLanguage) Print() string {
