XCore for GO v2
=============================

# Please use xcore/v2

# The version 1 is obsolete.

[![Go Report Card](https://goreportcard.com/badge/github.com/webability-go/xcore)](https://goreportcard.com/report/github.com/webability-go/xcore)
[![GoDoc](https://godoc.org/github.com/webability-go/xcore/v2?status.png)](https://godoc.org/github.com/webability-go/xcore/v2)
[![GolangCI](https://golangci.com/badges/github.com/webability-go/xcore.svg)](https://golangci.com)

# Please use xcore/v2

# The version 1 is obsolete.

import "github.com/webability-go/xcore/v2"

The XCore package is used to build basic object for programmation. for the WebAbility compatility code
For GO, the actual existing code includes:
- XCache: Application Memory Caches, thread safe.
- XDataset: Basic nested data structures for any purpose (template injection, configuration files, database records, etc) Support thread safe operations on thread safe structures (XDatasetTS and XDatasetCollectionTS)
- XLanguage: language dependent text tables, thread safe
- XTemplate: template system with meta language, thread safe cloning

Manuals are available on godoc.org [![GoDoc](https://godoc.org/github.com/webability-go/xcore/v2?status.png)](https://godoc.org/github.com/webability-go/xcore/v2)

# Please use xcore/v2

# The version 1 is obsolete.

Version Changes Control
=======================

v1.1.1 - 2022-09-02
-----------------------
- Bug corrected on XDataset.GetString() and XDatasetCollection.GetDataString(). 
  If the value is NIL int the dataset, it returns now "" and not "<nil>"

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
