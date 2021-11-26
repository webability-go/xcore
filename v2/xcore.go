// Copyright Philippe Thomassigny 2004-2020.
// Use of this source code is governed by a MIT licence.
// license that can be found in the LICENSE file.

// Package xcore is a set of basic objects for programation (XCache for caches, XDataset for data sets, XLanguage for languages and XTemplate for templates).
// For GO, the actual existing code includes:
//
// - XCache: Application Memory Caches for any purpose, with time control and quantity control of object in the cache and also check changes against original source. It is a thread safe cache.
//
// - XDataset: Basic nested data structures for any purpose (template injection, configuration files, database records, etc).
//
// - XLanguage: language dependent text tables for internationalization of code. The sources can be text or XML file definitions.
//
// - XTemplate: template system with meta language to create complex documents (compatible with any text language, HTML, CSS, JS, PDF, XML, etc), heavily used on CMS systems and others.
//
// It is already used on sites that serve more than 60 million pages a month (500 pages per second on pike hour) and can be used on multithreading environment safely.
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
// 2. Fill in the cache:
//
// Once you have declared the cache, you can fill it with anything you want. The main cache object is an interface{}
// so you can put here anything you need, from simple variables to complex structures. You need to use the Set function:
// Note the ID is always a string, so convert a database key to string if needed.
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
// 4. To maintain the cache:
//
// You may need Del function, to delete a specific entry (maybe because you deleted the record in database).
// You may also need Clean function to deletes a percentage of the cache, or Flush to deletes it all.
// The Verify function is used to check cache entries against their sources through the Validator function.
// Be very careful, if the cache is big or the Validator function is complex (maybe ask for a remote server information),
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
// 5. How to use Verify Function:
//
// This function is recommended when the source is local and fast to check (for instance a language file or a template file).
// When the source is distant (other cluster database, any rpc source on another network, integration of many parts, etc), it is more recommended to create a
// function that will delete the cache when needed (on demand cache change).
//
// The validator function is a func(id, time.Time) bool function. The first parameter is the ID entry in the cache, the second parameter the time of the entry was created.
// The validator function returns true is the cache is still valid, or false if it needs to be invalidated.
//
//  var myfiles = xcore.NewXCache("myfiles", 50, 0)
//  myfiles.Validator = FileValidator
//
//  // FileValidator verify the file source. In this case, the ID is directly the filename full path
//  func FileValidator(key string, otime time.Time) bool {
//
//	  fi, err := os.Stat(key)
//	  if err != nil {
//		  // Does not exists anymore, invalid
//		  return false
//	  }
//	  mtime := fi.ModTime()
//	  if mtime.After(otime) {
//		  // file is newer, invalid
//	  	return false
//  	}
//	  // All ok, valid
//	  return true
//  }
//
//
// The XCache is thread safe.
// The cache can be limited in quantity of entries and timeout for data. The cache is automanaged (for invalid expired data) and can be cleaned partially or totally manually.
//
//
// XLanguage
//
// The XLanguage table of text entries can be loaded from XML file, XML string or normal text file or string.
// It is used to keep a table of id=value set of entries in any languages you need, so it is easy to switch between XLanguage instance based on the required language needed.
// Obviously, any XLanguage you load in any language should have the same id entries translated, for the same use.
// The XLanguage object is thread safe
//
// 1. loading:
//
// You can load any file or XML string directly into the object.
//
// 1.1 The XML Format is:
//
//  <?xml version="1.0" encoding="UTF-8"?>
//  <language id="NAMEOFTABLE" lang="LG">
//    <entry id="ENTRYNAME">ENTRYVALUE</entry>
//    <entry id="ENTRYNAME">ENTRYVALUE</entry>
//    etc.
//  </language>
//
// NAMEOFTABLE is the name of your table entry, for example "loginform", "user_report", etc.
//
// LG is the ISO-3369 2 letters language ID, for example "es" for spanish, "en" for english, "fr" for french, etc.
//
// ENTRYNAME is the ID of the entry, for example "greating", "yourname", "submitbutton".
//
// ENTRYVALUE is the text for your entry, for example "Hello", "You are:", "Save" if your table is in english.
//
// 1.2 The flat text format is:
//
//  ENTRYNAME=ENTRYVALUE
//  ENTRYNAME=ENTRYVALUE
//  etc.
//
// ENTRYNAME is the ID of the entry, for example "greating", "yourname", "submitbutton".
//
// ENTRYVALUE is the text for your entry, for example "Hello", "You are:", "Save" if your table is in english.
//
// There is no name of table or language in this format (you "know" what you are loading).
//
// The advantage to use XML format is to have more control over your language, and eventyally add attributes into your entries,
// for instance you may add attributes translated="yes/no", verified="yes/no", and any other data that your system could insert.
// The XLanguage will ignore those attributes loading the table.
//
// 2. creation:
//
// To create a new XLanguage empty structure:
//
//  lang := NewXLanguage(id, language)
//
// There are 4 functions to create the language from a file or string, flat text or XML text:
//
//  langfromxmlfile := NewXLanguageFromXMLFile("path/to/filename")
//  langfromxmlstring := NewXLanguageFromXMLString("<xml>...")
//  langfromtextfile := NewXLanguageFromFile("path/to/file")
//  langfromtextstring := NewXLanguageFromString("entry=data\n...")
//
// Then you can use the set of basic access functions:
//
// SetName/SetLanguage functions are used to set the table name and language of the object (generally to build an object from scratch).
// GetName/GetLanguage functions are used to get the table name and language of the object (generally when you load it from some source).
// Set/Get/Del functions are used to add or modify a new entry, read an entry, or deletes an entry in the object.
//
// XDataSet
//
// 1. Overview:
//
// The XDataSet is a set of interfaces and basic classes ready-to-use to build a standard set of data optionally nested and hierarchical, that can be used for any purpose:
//
// - Keep complex data in memory.
//
// - Create JSON structures.
//
// - Inject data into templates.
//
// - Interchange database data (records set and record).
//
// You can store into it generic supported data, as well as any complex interface structures:
//
// - Int
//
// - Float
//
// - String
//
// - Time
//
// - Bool
//
// - []Int
//
// - []Float
//
// - []Time
//
// - []Bool
//
// - XDataSetDef (anything extended with this interface)
//
// - []String
//
// - Anything else ( interface{} )
//
// - XDataSetCollectionDef (anything extended with this interface)
//
// The generic supported data comes with a set of functions to get/set those data directly into the XDataset.
//
// Example:
//
//  import "github.com/webability-go/xcore"
//
//  data := xcore.XDataset{}
//  data["data1"] = "DATA1"
//  data["data2"] = "DATA1"
//  sm := xcore.XDataset{}
//  sm["data31"] = "DATA31"
//  data["data3"] = sm
//  data["data4"] = 123
//  data["data5"] = 123.432
//  data["data6"] = true
//  data["data7"] = func() string { return "ABC" }
//
//  d8_r1 := &xcore.XDataset{}
//  d8_r1.Set("data81", "rec 1: Entry 8-1")
//  d8_r1.Set("data82", "rec 1: Entry 8-2")
//
//  d8_r2 := &xcore.XDataset{}
//  d8_r2.Set("data81", "rec 2: Entry 8-1")
//  d8_r2.Set("data82", "rec 2: Entry 8-2")
//  d8_r2.Set("data83", "rec 2: Entry 8-3")
//
//  d8_r3 := &xcore.XDataset{}
//  d8_r3.Set("data81", "rec 3: Entry 8-1")
//  d8_r3.Set("data82", "rec 3: Entry 8-2")
//
//  d := xcore.XDatasetCollection{}
//  d.Push(d8_r1)
//  d.Push(d8_r2)
//  d.Push(d8_r3)
//
//  data["data8"] = &d
//  data["data9"] = "I exist"
//
// Note that all references to XDataset and XDatasetCollection are pointers, always (to be able to modify the values of them).
//
//
// 2. XDatasetDef interface:
//
// It is the interface to describe a simple set of data mapped as "name": value, where value can be of any type.
//
// The interface implements a good amount of basic methods to get the value on various format such as GetString("name"), GetInt("name"), etc (see below).
//
// If the value is another type as asked, the method should contert it if possible. For instance "key":123 required through GetString("key") should return "123".
//
// The XDataset type is a simple map[string]interface{} with all the implemented methods and should be enough to use for almost all required cases.
//
// However, you can build any complex structure that extends the interface and implements all the required functions to stay compatible with the XDatasetDef.
//
// 3. XDatasetCollectionDef Interface:
//
// This is the interface used to extend any type of data as a Collection, i-e an array of XDatasetDef. This is a slice of any XDatasetDef compatible data.
//
// The interface implements some methods to work on array structure such as Push, Pop, Shift, Unshift and some methods to search data into the array.
//
// The XDatasetCollection type is a simple []DatasetDef with all the implemented methods and should be enough to use for almost all required cases.
//
// XDataSetTS
//
// 1. Overview:
//
// The XDataSetTS is a DatasetDef structure, thread safe.
// It is build on the XDataset with the same properties, but is thread safe to protect Read/Write accesses from different thread.
//
// Example:
//
//  import "github.com/webability-go/xcore/v2"
//
//  data := &xcore.XDatasetTS{} // data is a XDatasetDef
//  data.Set("data1", "DATA1")
//  data.Set("newkey", 123.45)
//
// You may also build a XDatasetTS to encapsulate a XDatasetDef that is not thread safe, to use it safely
//
//  import "github.com/webability-go/xcore/v2"
//
//  datanots := xcore.NewXDataset()
//  datats := xcore.NewXDatasetTS(datanots)
//
// Note that all references to XDatasetTS are pointers, always (to be able to modify the values of them).
//
// The DatasetTS meet the XDatasetDef interface
//
//
// XTemplate
//
// 1. Overview:
//
// This is a class to compile and keep a Template that can be injected with an XDataSet structure of data, with a metalanguage to inject the data.
//
// The metalanguage is extremely simple and is made to be useful and **really** separate programation from template code (not like other many generic template systems that just mix code and data).
//
// A template is a set of HTML/XML (or any other language) string with a meta language to inject variables and build a final string.
//
// The XCore XTemplate system is based on the injection of parameters, language translation strings and data fields directly into the HTML (Or any other language you need) template.
//
// The HTML itself (or any other language) is a text code not directly used by the template system, but used to dress the data you want to represent in your preferred language.
//
// The variables to inject must be into a XDataSet structure or into a structure extended from XDataSetDef interface.
//
// The injection of data is based on a XDataSet structure of values that can be nested into another XDataSet and XDataSetConnection and so on.
//
// The template compiler recognize nested arrays to automatically make loops on the information.
//
// Templates are made to store reusable HTML code, and overall easily changeable by people that do not know how to write programs.
//
// A template can be as simple as a single character (no variables to inject) to a very complex nested, conditional and loops sub-templates.
//
//  Hello world!
//
// Yes. this is a template, but a very simple one without need to inject any data.
//
// Let's go more complex:
//
// Having an array of data, we want to paint it beautifull:
//
//  { "clientname": "Fred",
//    "hobbies": [
//       { "name": "Football" },
//       { "name": "Ping-pong" },
//       { "name": "Swimming" },
//       { "name": "Videogames" }
//    ]
//  }
//
// We can create a template to inject this data into it:
//
//  %-- This is a comment. It will not appear in the final code. --%
//  Let's put your name here: {{clientname}}<br />
//  And lets put your hobbies here:<br />
//  @@hobbies:hobby@@     %-- note the 1rst id is the entry into the data to inject and the second one is the name of the sub-template to use --%
//
//  %-- And you need the template for each hobby:--%
//  [[hobby]]
//  I love {{name}}<br />
//  [[]]
//
//
// 2. Create and use XTemplateData:
//
// In sight to create and use templates, you have all those possible options to use:
//
// Creates the XTemplate from a string or a file or any other source:
//
//  package main
//
//  import (
//    "fmt"
//    "github.com/webability-go/xcore"
//  )
//
//  func main() {
//    tmpl, _ := xcore.NewXTemplateFromString(`
//  %-- This is a comment. It will not appear in the final code. --%
//  Let's put your name here: {{clientname}}<br />
//  And lets put your hobbies here:<br />
//  @@hobbies:hobby@@     %-- note the 1rst id is the entry into the data to inject and the second one is the name of the sub-template to use --%
//
//  %-- And you need the template for each hobby:--%
//  [[hobby]]
//  I love {{name}}<br />
//  [[]]
//  `)
//    // The creation of the data is obviously tedious here, in real life it should come from a JSON, a Database, etc
//    data := xcore.XDataset{
//      "clientname": "Fred",
//      "hobbies": &XDatasetCollection{
//        &XDataset{"name":"Football"},
//        &XDataset{"name":"Ping-pong"},
//        &XDataset{"name":"Swimming"},
//        &XDataset{"name":"Videogames"},
//      },
//    }
//
//    fmt.Println(tmpl.Execute(&data)
//  }
//
// Clone the XTemplate:
//
//  xtemplate := xcore.NewXTemplate()
//  xtemplatecloned := xtemplate.Clone()
//
//
// 3. Metalanguage Reference:
//
// 3.1 Comments: %-- and --%
//
// You may use comments into your template.
// The comments will be discarded immediately at the compilation of the template and do not interfere with the rest of your code.
//
// Example:
//
//  %-- This is a comment. It will not appear in the final code. --%
//
//  %--
//  This subtemplate will not be compiled, usable or even visible since it is into a comment
//  [[templateid]]
//  Anything here
//  [[]]
//  --%
//
//
// 3.2 Nested Templates: [[...]] and [[]]
//
// You can define new nested templates into your main template
// A nested template is defined by:
//
//  [[templateid]]
//  your nested template here
//  [[]]
//
// The templteid is any combination of lowers letters only (a-z), numbers (0-9), and 3 special chars: . (point) - (dash) and _ (underline).
//
// The template is closed with [[]].
//
// There is no limits into nesting templates.
//
// Any nested template will inheritate all the father elements and can use father elements too.
//
// To call a sub-template, you need to use &&templateid&& syntax (described below in this document).
//
// Example:
//
//  &&header&&
//  Welcome to my page
//  &&footer&&
//
//  [[header]]
//  <hr />
//  [[]]
//
//  [[footer]]
//  <hr />
//  &&copyright&&
//
//    [[copyright]]
//      © 2012 Ing. Philippe Thomassigny, a project of the WebAbility® Network.
//    [[]]
//  [[]]
//
// You may use more than one id into the same template to avoid repetition of the same code.
// The different id's are separated with a pipe |
//
//  [[looptemplate_first|looptemplate_last|looptemplate_odd]] %-- First element, last element and all odd elements will be red --%
//    <div style="color: red;">{{data}}</div>
//  [[]]
//  [[looptemplate]] %-- All the other elements will be blue --%
//    <div style="color: blue;">{{data}}</div>
//  [[]]
//
// Important note:
// A template will be visible only on the same level of its declaration. For example, if you put a subtemplate "b" into a subtemplate "a", it will not be visible by &&b&& from the top level, but only into the subtemplate "a".
//
//  &&header&&
//  Welcome to my page
//  &&footer&&
//  &&copyright&& %-- WILL NOT WORK, into a sub-template --%
//
//  [[header]]
//  <hr />
//  [[]]
//
//  [[footer]]
//  <hr />
//  &&copyright&& %-- WILL WORK, same level --%
//
//    [[copyright]]
//      © 2012 Ing. Philippe Thomassigny, a project of the WebAbility® Network.
//    [[]]
//  [[]]
//
//
// 3.3 Simple Elements: ##...## and {{...}}
//
// There are 2 types of simple elements. Language elements and Data injector elements (also called field elements).
//
// We "logically" define the 2 type of elements. The separation is only for human logic and template filling, however the language information can perfectly fit into the data to inject (and not use ## entries).
//
//
// 3.3.1 Languages elements: ##entry##
//
// All the languages elements should have the format: ##entry##.
//
// A language entry is generally anything written into your code or page that does not come from a database, and should adapt to the language of the client visiting your site.
//
// Using the languages elements may depend on the internationalization of your page.
//
// If your page is going to be in a single language forever, you really dont need to use languages entries.
//
// The language elements generally carry titles, menu options, tables headers etc.
//
// The language entries are set into the "#" entry of the main template XDataset to inject, and is a XLanguage table.
//
// Example:
//
//  <div style="background-color: blue;">
//  ##welcome##<br />
//  You may use the same parameter as many time you wish.<br />
//
//  <span onclick="alert('##hello##');" class="button">##clickme##!</span>
//  <span onclick="alert('##helloagain##');" class="button">##clickme## ##again##!</span>
//
//  </div>
//
// With data to inject:
//
//  {
//    "#": {
//      "welcome": "Bienvenue",
//      "clickme": "Clique sur moi",
//      "again": "de nouveau"
//    }
//  }
//
//
// 3.3.2 Field elements: {{fieldname}}
//
// Fields values should have the format: {{fieldname}}.
//
// Your fields source can be a database or any other preferred repository data source.
//
// Example:
//
//  <div style="background-color: blue;">
//  Welcome to my site<br />
//  You may use the same parameter as many time you wish.<br />
//
//  Today's temperature is {{degres}} celsius<br />
//  or {{fdegres}} farenheit<br />
//
//  Is {{degres}} degres too cold ? Buy a pullover!<br />
//
//  </div>
//
// You can access an element with its path into the data set to inject separating each field level with a > (greater than).
//
//  {{hobbies>1>name}}
//
// This will take the name of the second hobby in the dataset defined upper. (collections are 0 indexed).
//
// The 1 denotes the second record of the hobbies XDatasetCollection.
//
// If the field is not found, it will be replaced with an empty string.
//
// Tecnically your field names can be any string in the dataset. However do not use { } or > into the names of your fields or the XTemplate may not use them correctly.
//
// We recommend to use lowercase names with numbers and ._- Accents and UTF8 symbols are also welcome.
//
// 3.3.3 Scope:
//
// When you use an id to point a value, the template will first search into the available ids of the local level.
// If no id is found, the it will search into the upper levers if any, and so on.
//
// Example:
//
//  {
//    "detail": {
//      "data1": {
//        "data2": {
//          "key1": {"appname": "Nested App", "name": "Juan", "status": 1},
//          "key2": {"name": "José", "status": 2},
//          "appname" => "DomCore"
//        }
//      }
//    }
//  }
//
// At the level of 'data2', using {{appname}} will get back 'DomCore'.
//
// At the level of 'key1', using {{appname}} will get back 'Nested App'.
//
// At the level of 'key2', using {{appname}} will get back 'DomCore'.
//
// At the level of root, 'data1' or 'detail', using {{appname}} will get back an empty string.
//
//
// 3.3.4 Path access: id>id>id>id
//
// At any level into the data array, you can access any entry into the subset array.
//
// For instance, taking the previous array of data to inject,
// let's suppose we are into a nested meta elements at the 'data1' level. You may want to access directly the 'Juan' entry. The path will be:
//
//  {{data2>key1>name}}
//
// The José's status value from the root will be:
//
//  {{detail>data1>data2>key2>status}}
//
//
//
// 3.4 Meta Elements
//
// They consist into an injection of a XDataset, called the "data to inject", into the template.
//
// The meta language is directly applied on the structure of the data array.
//
// The data to inject is a nested set of variables and values with the structure you want (there is no specific construction rules).
//
// You can inject nearly anything into a template meta elements.
//
// Example of a data array to inject:
//
//    data := xcore.XDataset{
//      "clientname": "Fred",
//      "clientpicture": "face.jpg",
//      "hobbies": &XDatasetCollection{
//        &XDataset{"name":"Football","sport":"yes"},
//        &XDataset{"name":"Ping-pong","sport":"yes"},
//        &XDataset{"name":"Swimming","sport":"yes"},
//        &XDataset{"name":"Videogames","sport":"no"},
//      },
//      "preferredhobby": &XDataset{
//        "name":"Baseball",
//        "sport":"yes",
//      },
//      "metadata": &XDataset{
//        "preferred-color": "blue",
//        "Salary": 3568.65,
//        "hiredate": time.Time("2020-01-01T12:00:00"),
//      },
//    }
//
// You can access directly any data into the array with its relative path (relative to the level you are when the metaelements are applied, see below).
//
// There are 4 structured meta elements in the XTemplate templates to use the data to inject:
//
// Reference, Loops, Condition and Debug.
//
// The structure of the meta elements in the template must follow the structure of the data to inject.
//
//
// 3.4.1 References to another template: &&order&&
//
// 3.4.1.1 When order is a single id (characters a-z0-9.-_), it will make a call to a sub template with the same set of data and replace the &&...&& with the result.
// The level in the data set is not changed.
//
// Example based on previous array of Fred's data:
//
//  &&header&&
//  &&body&&
//
//  [[header]]
//  Sports shop<hr />
//  [[]]
//
//  [[body]]
//  {{clientname}} Data:
//  <img src="{{clientpicture}}" title="{{clientname}}" />
//  [[]]
//
//
// 3.4.1.2 When order contains 2 parameters separated by a semicolumn :, then second parameter is used to change the level of the data of array, with the subset with this id.
// The level in the data set is changed to this sub set.
//
// Example based on previous array of Fred's data:
//
//  &&header&&
//  &&body:metadata&&
//
//  [[header]]
//  Sports shop<hr />
//  [[]]
//
//  [[body]]
//  {{clientname}} Data:  %-- Taken from the root data --%
//  Salary: {{salary}}<br /> %-- taken from the metadata subset --%
//  Hire date: {{hiredate}}<br /> %-- taken from the metadata subset--%
//  [[]]
//
// 3.4.1.3 When order contains 3 parameters separated by a semicolumn :, the second and third parameters are used to search the name of the new template based on the data fields to inject.
//
// This is an indirect access to the template. The name of the subtemplate is build with parameter3 as prefix and the content of parameter2 value.
// The third parameter must be empty.
//
//  &&header&&
//  &&body:preferredhobby&&
//
//  [[header]]
//  Sports shop<hr />
//  [[]]
//
//  [[body]]
//  {{clientname}} Preferred hobby:
//  &&:sport:sport.&&  %-- will build sport_ + [yes/no] contained into the sport field. Be sure you have a template for each value ! --%
//
//  [[sport.yes]]{{name}} - It's a sport, sell him things![[]]
//  [[sport.no]]{{name}} - It's not a sport, recommend him next store.[[]]
//  [[sport]]{{name}} - We do not know that it is.[[]]
//  [[]]
//
//
// 3.4.2 Loops: @@order@@
//
// 3.4.2.1 Overview
//
// This meta element will loop over each itterance of the set of data and concatenate each created template in the same order. You need to declare a sub template for this element.
//
// You may aso declare derivated sub templates for the different possible cases of the loop:
// For instance, If your main subtemplate for your look is called "hobby", you may need a different template for the first element, last element, Nth element, Element with a value "no" in the sport field, etc.
//
// The supported postfixes are:
//
// When the array to iterate is empty:
//
// - .none (for example "There is no hobby")
//
// When the array contains elements, it will search in order, the following template and use the first found:
//
// - templateid.key.[value]  value is the key of the vector line. If the collection has a named key (string) or is a direct array (0, 1, 2...)
//
// - templateid.first if it is the first element of the array set (new from v1.01.11)
//
// - templateid.last if it is the first element of the array set (new from v1.01.11)
//
// - templateid.even if the line number is even
//
// - templateid in all other cases (odd is contained here if even is defined)
//
//
// 3.4.2.2 When order is a single id (characters a-z0-9.-_), it will make a call to the sub template id with the same subset of data with the same id and replace the @@...@@ for each itterance of the data with the result.
//
// Example based on previous array of Fred's data:
//
//  &&header&&
//  &&body&&
//
//  [[header]]
//  Sports shop<hr />
//  [[]]
//
//  [[body]]
//  {{clientname}} Data:
//  <img src="{{clientpicture}}" title="{{clientname}}" />
//  Main hobby: {{preferredhobby>name}}<br />
//  Other hobbies:<br />
//  @@hobbies@@
//  [[hobbies.none]]There is no hobby<br />[[]]
//  [[hobbies]]{{name}}<br />[[]]
//  [[]]
//
// 3.4.2.3 When order contains 2 parameters separated by a semicolumn :, then first parameter is used to change the level of the data of array, with the subset with this id, and the second one for the template to use.
//
// Example based on previous array of Fred's data:
//
//  &&header&&
//  &&body&&
//
//  [[header]]
//  Sports shop<hr />
//  [[]]
//
//  [[body]]
//  {{clientname}} Data:
//  <img src="{{clientpicture}}" title="{{clientname}}" />
//  Main hobby: {{preferredhobby>name}}<br />
//  Other hobbies:<br />
//  @@hobbies:hobby@@ %-- will iterate over hobbies in the data, but with hobby sub template --%
//  [[hobby.none]]There is no hobby<br />[[]]
//  [[hobby.key.1]]<span style="color: red;">{{name}}<span><br /><hr>[[]] %-- Paint the second line red (0 indexed) --%
//  [[hobby]]{{name}}<br />[[]]
//  [[]]
//
//
// 3.4.3 Conditional: ??order??
//
// Makes a call to a subtemplate only if the field exists and have a value.
// This is very userfull to call a sub template for instance when an image or a video is set.
//
// When the condition is not met, it will search for the [id].none template.
// The conditional element does not change the level in the data set.
//
// 3.4.3.1 When order is a single id (characters a-z0-9.-_), it will make a call to the sub template id with the same field in the data and replace the ??...?? with the corresponding template
//
// Example based on previous array of Fred's data:
//
//  &&header&&
//  &&body&&
//
//  [[header]]
//   Sports shop<hr />
//  [[]]
//
//  [[body]]
//   {{clientname}} Data:
//   ??clientpicture??
//   [[clientpicture]]<img src="{{clientpicture}}" title="{{clientname}}" />[[]]
//   [[clientpicture.none]]There is no photo<br />[[]]
//  [[]]
//
// 3.4.3.2 When order contains 2 parameters separated by a semicolumn :, then second parameter is used to change the level of the data of array, with the subset with this id.
//
// Example based on previous array of Fred's data:
//
//  &&header&&
//  &&body&&
//
//  [[header]]
//   Sports shop<hr />
//  [[]]
//
//  [[body]]
//   {{clientname}} Data:
//   ??clientpicture:photo??
//   [[photo]]<img src="{{clientpicture}}" title="{{clientname}}" />[[]]
//   [[photo.none]]There is no photo<br />[[]]
//  [[]]
//
// If the asked field is a catalog, true/false, numbered, you may also use .[value] subtemplates
//
//  &&header&&
//  &&body&&
//
//  [[header]]
//   Sports shop<hr />
//  [[]]
//
//  [[body]]
//   {{clientname}} Data:
//   ??preferredhobby>sport:preferredhobby??
//   [[preferredhobby.yes]]{{preferredhobby>name}}<br />[[]]
//   [[preferredhobby|preferredhobby.no|preferredhobby.none]]There is no preferred sport<br />[[]]
//  [[]]
//
//
// 3.5 Debug Tools: !!order!!
//
// There are two keywords to dump the content of the data set.
// This is very useful when you dont know the code that calls the template, don't remember some values, or for debug facilities.
//
// 3.5.1 !!dump!!
//
// Will show the totality of the data set, with ids and values.
//
// 3.5.1 !!list!!
//
// Will show only the tree of parameters, values are not shown.
//
//
package xcore

// VERSION is the used version nombre of the XCore library.
const VERSION = "2.0.9"

// LOG is the flag to activate logging on the library.
// if LOG is set to TRUE, LOG indicates to the XCore libraries to log a trace of functions called, with most important parameters.
// LOG can be set to true or false dynamically to trace only parts of code on demand.
var LOG = false
