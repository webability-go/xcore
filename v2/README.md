XCore v2 for GO
=============================

[![Go Report Card](https://goreportcard.com/badge/github.com/webability-go/xcore)](https://goreportcard.com/report/github.com/webability-go/xcore)
[![GoDoc](https://godoc.org/github.com/webability-go/xcore/v2?status.png)](https://godoc.org/github.com/webability-go/xcore/v2)
[![GolangCI](https://golangci.com/badges/github.com/webability-go/xcore.svg)](https://golangci.com)

Minimum version of GO: 1.17 (for time.Time compatibility)

The XCore package is used to build basic object for programmation. for the WebAbility compatility code
For GO, the actual existing code includes:
- XCache: Application Memory Caches, thread safe.
- XDataset: Basic nested data structures for any purpose (template injection, configuration files, database records, etc) Support thread safe operations on thread safe structures (XDatasetTS and XDatasetCollectionTS)
- XLanguage: language dependent text tables, thread safe
- XTemplate: template system with meta language, thread safe cloning

Manuals are available on godoc.org [![GoDoc](https://godoc.org/github.com/webability-go/xcore/v2?status.png)](https://godoc.org/github.com/webability-go/xcore/v2)


TO DO, maybe:
=============
- XDataset.Set should accept path too > > >
- Get*Collection should convert types too
- XTemplate must concatenate strings after compilation
- Implements functions as data entry for template Execute (simple data or loop functions, can get backs anything, creates an interface)
Some improvements to check, later:
- Adds mutex on XTemplate ?? (they should be used locally on every thread, or not ??), maybe adds a flag "thread safe" ?
- XCache: activate persistant cache too (shared memory) ????? maybe not for go itself, but for instance to talk with other memory data used by other languages and apps, or to not loose the caches if the app is restarted.


Version Changes Control
=======================

v2.2.3 - 2025-06-09
-----------------------
- Corrected a bug into xtemplate that made the conditional for sub templates with an array of dataset not working

v2.2.2 - 2023-10-12
-----------------------
- Added a security on the sub template .none for ?? meta language into XTemplate, to not try to use an inexistant template and throw a panic error.

v2.2.1 - 2023-09-29
-----------------------
- Added the missing sub template .none for ?? meta language into XTemplate.

v2.2.0 - 2023-08-31
-----------------------
- Added the parameter status to xlanguage XML and function to get/set the parameter.
- Added the function GetXML() to marshal the structure to an XML file.

v2.1.7 - 2023-06-15
-----------------------
- Added the missing counter for @@ meta language. Using {{.counter}} into the loop template, you can add the number of the dataset, 1-based.

v2.1.6 - 2023-05-25
-----------------------
- Bug corrected into GetCollection(id) function, if the dataset is not an XDatasetCollection, it should return nil, not a panic error

v2.1.5 - 2023-05-17
-----------------------
- Bump golang.org/x/text from 0.3.5 to 0.3.8 for security upgrade

v2.1.4 - 2023-05-17
-----------------------
- Bug corrected on @@ loops subtemplates, the .none template was not reach when the array is set but empty

v2.1.3 - 2022-09-02
-----------------------
- Bug corrected on XDataset.GetString() and XDatasetCollection.GetDataString(). 
  If the value is NIL int the dataset, it returns now "" and not "<nil>"

v2.1.2 - 2022-03-02
-----------------------
- XTemplate: Added = to metalanguage string tags to resolve also the paths (bug corrected).

v2.1.1 - 2022-03-01
-----------------------
- XTemplate: Added > to metalanguage string tags to resolve also the paths (bug corrected).

v2.1.0 - 2022-02-27
-----------------------
- XLanguage: bug corrected on unlock of stringload and loadFromFile (was blocking the system)
- XTemplate: The metalanguage keywords are now only recognized if they match authorized characters (for instance &&keyword&&), to avoid bugs in JS with && and || and !!.
- Print functions of time.Time corrected (as in go 1.17, the print format changes) into the *test.go test functions

v2.0.9 - 2021-11-25
-----------------------
- XCache modified to defer mutex unlocks instead of directly unlock into the code, to avoid dead locks in case of thread panic and crashes.
- XLanguage modified to defer mutex unlocks instead of directly unlock into the code, to avoid dead locks in case of thread panic and crashes.

v2.0.8 - 2021-05-18
-----------------------
- XTemplate is now clonable: newtemplate := template.Clone()

v2.0.7 - 2021-03-23
-----------------------
- []interface{} added to NewXDataset to be decoded as []map[string]interface{} because of decoded JSON


v2.0.5 - 2021-03-23
-----------------------
- function NewXDataset(data) and NewXDatasetCollection(data) added, to build XDatasets based on a classic map[string]interface{}
- XDatasetTS added to the manual
- XLanguage is now thread safe
- XTemplate: Error corrected in string() function (was saying XLanguage)
- Each cache entry is now able to manage its own TLL if set. New function SetTTL(id, duration)

v2.0.4 - 2020-04-13
-----------------------
- XLanguage: Added GetEntries() func (not thread safe yet)

v2.0.3 - 2020-04-08
-----------------------
- XLanguage: Error corrected on loadXMLString: the data was not loading correctly into the XLanguage object.

v2.0.1, v2.0.2 - 2020-03-29
-----------------------
- Version adjustment for github and go modules v2

v2.0.0 - 2020-03-29
-----------------------
- xdataset.go now as a coverage of 100% with xdataset_test.go
- XCache now uses R/W mutex
- New interfaces.go file to keep all the interfaces in it (XDatasetDef, XDatasetCollectionDef)
- New xdatasetts.go for thread safe dataset
- New xdatasetts_test.go for thread safe dataset tests
- New xdatasetcollection.go for the collection of dataset (separation from xdataset.go)
- New xdatasetcollection_test.go for collection tests
- New xdatasetcollectionts.go for thread safe dataset
- New xdatasetcollectionts_test.go for thread safe datasetcollection tests
- XLanguage is now thread safe with R/W mutexes

v1.1.0 - 2020-03-01
-----------------------
- Modularization of XCore
- XLanguage tests and examples are now conform to Go test units
- Implementation of XLanguage.String and XLanguage.GoString, removed Print
- XCache tests and examples are now conform to Go test units
- XDataset tests and examples are now conform to Go test units
- Implementation of XDataset.String and XDataset.GoString, removed Print
- Implementation of XDatasetCollection.String and XDatasetCollection.GoString, removed Print
- XTemplate tests and examples are now conform to Go test units
- Implementation of XTemplate.String and XTemplate.GoString, removed Print

v1.0.1 - 2020-02-10
-----------------------
- Documentation corrections
- Bug on String() and GoString() corrected

v1.0.0 - 2020-02-09
-----------------------
- Version leveling
- Documentation corrections
- Change functions Stringify() by String() and GoString() for language compatibility
- Tests functions enhanced

v0.3.1 - 2020-02-09
-----------------------
- XDatasetDef.Get must accept a path as key (id>id>id)
- XTemplates now resolve {{ fields with path id>id>id
- XTemplates now resolve @@ metalanguage with 1 and 2 Parameters
- XTemplates now resolve && metalanguage with 1,2 and 3 Parameters
- XTemplates now resolve ?? metalanguage with 1, and 2 Parameters
- XTemplates now resolve !! debug orders
- XTemplates now implements sub templates derivation (.none .first .last .(number) )
- Manuals for XCache, XLanguage and XTemplate written with reference of the metalanguage
- Examples for dataset and xtemplate added (working version)
- XDataset and XDatasetCollection .Stringify now prints also field names.

v0.3.0 - 2020-02-06
-----------------------
- The properties of XTemplateParam are now public so the full structure can be used to build other type of code based on the XTemplate rules
- The subtemplates IDs must be lowers, numbers and . - _ in sight of integration with other systems that can mix tags [[]] within the code

v0.2.3 - 2020-01-23
-----------------------
- Corrected a bug to avoid null pointer panic error if the array of data for XTemplate.Execute function is nil

v0.2.2 - 2019-12-21
-----------------------
- XLanguage now support golang x/text/language instead of direct iso 2 charater language
- godoc manuals for xlanguage, xdataset and xtemplate prepared

v0.2.1 - 2019-12-13
-----------------------
- XCache manual enhanced with examples

v0.2.0 - 2019-12-06
-----------------------
- XCache Code simplified to expose XCache definition as public, remove not usefull funcion (Get*)
- XCache 0.2.0 is not compatible with XCache 0.1.* , you may need to change your code
- Added more conversions between int-float-bool in XDataset.Get*

v0.1.2 - 2019-12-05
-----------------------
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
