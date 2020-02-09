@UTF-8

[![Go Report Card](https://goreportcard.com/badge/github.com/webability-go/xcore)](https://goreportcard.com/report/github.com/webability-go/xcore)
[![GoDoc](https://godoc.org/github.com/webability-go/xcore?status.png)](https://godoc.org/github.com/webability-go/xcore)
[![GolangCI](https://golangci.com/badges/github.com/webability-go/xcore.svg)](https://golangci.com)

XCore for GO v0
=============================

The XCore package is used to build basic object for programmtion. for the WebAbility compatility code
For GO, the actual existing code includes:
- XCache: Application Memory Caches
- XDataset: Basic nested data structures for any purpose (template injection, configuration files, database records, etc)
- XLanguage: language dependent text tables
- XTemplate: template system with meta language

Manuals are available on godoc.org [![GoDoc](https://godoc.org/github.com/webability-go/xcore?status.png)](https://godoc.org/github.com/webability-go/xcore)


TO DO:
======
- Implement Logging with import "log"
- XTemplate must concatenate strings after compilation
- Implements functions as data entry for template Execute (simple data or loop functions, can get backs anything, creates an interface)
- template.Print beautify, check stringify
- language.Print beautify, check stringify
- Some improvements to check, later:
XCache: activate persistant cache too (shared memory) ????? maybe not for go itself, but for instance to talk with other memory data used by other languages and apps, or to not loose the caches if the app is restarted.


Version Changes Control
=======================

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


XLanguage:

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
