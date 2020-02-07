// Copyright Philippe Thomassigny 2004-2020
// Use of this source code is governed by a MIT licence
// license that can be found in the LICENSE file.

// Package xcore is a set of basic objects for programation (XCache for caches, XDataset for data sets, XLanguage for languages and XTemplate for templates).
// For GO, the actual existing code includes:
//
// - XCache: Application Memory Caches for any purpose,
//
// - XDataset: Basic nested data structures for any purpose (template injection, configuration files, database records, etc),
//
// - XLanguage: language dependent text tables,
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
// The XCache is thread safe.
// The cache can be limited in quantity of entries and timeout for data. The cache is automanaged (for invalid expired data) and can be cleaned partially or totally manually.
//
//
// XLanguage
//
// The XLanguage table of text entries can be loaded from XML file, XML string or normal text file or string.
//
// 1. loading
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
// ENTRYNAME is the ID of the entry, for example "greating", "yourname", "submitbutton"
// ENTRYVALUE is the text for your entry, for example "Hello", "You are:", "Save" if your table is in english
//
// 1.2 The flat text format is:
//
//  ENTRYNAME=ENTRYVALUE
//  ENTRYNAME=ENTRYVALUE
//  etc.
//
// ENTRYNAME is the ID of the entry, for example "greating", "yourname", "submitbutton"
// ENTRYVALUE is the text for your entry, for example "Hello", "You are:", "Save" if your table is in english
//
// There is no name of table or language in this format (you "know" what you are loading)
//
// The advantage to use XML format is to have more control over your language, and eventyally add attributes into your entries,
// for instance translated="yes/no", verified="yes/no", and any other data that your system could insert
//
// 2. creation
//
// To create a new XLanguage empty structure:
//
//  lang := NewXLanguage()
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
// SetName/SetLanguage functions are used to set the table name and language of the object (generally to build an object from scratch)
// GetName/GetLanguage functions are used to get the table name and language of the object (generally when you load it from some source)
// Set/Get/Del functions are used to add or modify a new entry, read an entry, or deletes an entry in the object.
//
// XDataSet
//
// 1. Overview
//
// The XDataSet is a set of interfaces and basic classes ready-to-use to build a standard set of data optionally nested and hierarchical, that can be used for any purpose:
//
// - Keep complex data in memory
//
// - Create JSON structures
//
// - Inject data into templates
//
// - Interchange database data (records set and record)
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
// 2. XDatasetDef interface
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
//
// XTemplate
//
// 1. Overview
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
// Templates are made to store reusable HTML code, and overall easily changeable by **NON PROGRAMMING PEOPLE**.
//
// A template can be as simple as a single character (no variables to inject) to a very complex nested, conditional and loops sub-templates.
//
//  Hello world!
//
// Yes. this is a template, but a very simple one without need to inject any data.
//
// Let's go more complex:
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
// The data to inject could be:
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
// 2. Create and use XTemplateData
//
// In sight to create and use templates, you have all those possible options to use:
//
// Creates the XTemplat from a string or a file or any other source:
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
//
// 3. Metalanguage Reference
//
// ** Comments
//
// You may use comments into your template.
// The comments will be discarded immediately at the compilation of the template and not interfere with the rest of your code.
//
// Comments are defined by %-- and --%
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
/*
* Nested templates, to store many pieces of HTML
* Simple elements, to replace by values in the template. There are various types of simple elements:
** Parameters
** Language entries
** Fields of values
* Meta elements, to build a code based on a data array. There are various types of meta elements:
** Data access with ~~{{...}}~~, to show the value of a data into the data array
** Subtemplates access with ~~&&...&&~~, to call a subtemplate based on the value of an entry in the data array.
** Conditional access with ~~??...??~~, to show a piece of HTML based on the existence or value of an entry in the data array.
** Loops access with ~~@@...@@~~, to repeat a piece of HTML based on a table of values.
** Debug tools with ~~!!...!!~~, to show the data array.


```
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
```




++ Nested Templates

You can define new nested templates into your main template
A nested template is defined by:

```
<pre>
~~[[templateid]]~~
your nested template
~~[[]]~~
</pre>
```

The id is any combination of letters (a-z, A-Z, accents are welcome too), numbers (0-9), and 3 special chars: .-_

The old syntax also work and is deprecated. It will be definitively removed as of beginning of 2013.

```
<pre>
%%SUBTEMPLATE(templateid)%%
your nested template
%%ENDSUBTEMPLATE%%
</pre>
```

There is no limits into nesting templates.
Any nested template will inheritate all the father elements and can use father elements too.

Example:

```
<pre>&&header&&
Welcome to my page
&&footer&&

~~[[header]]~~
  &lt;hr />
~~[[]]~~

~~[[footer]]~~
  &lt;hr />
  &&amp;copyright&&

  ~~[[copyright]]~~
    © 2012 Ing. Philippe Thomassigny, a project of the WebAbility® Network.
  ~~[[]]~~

~~[[]]~~
</pre>
```

You may use more than one id into the same template to avoid repetition of the same code.
The different id's are separated with a pipe |

```
<pre>
~~[[templateid|anotherid|something.key|andmoreid]]~~
your nested template
~~[[]]~~
</pre>
```




++ Elements

The DomCore template system is based on a parameters replacement and a simple macrolanguage keywords.
Note: The syntax of the parameters, languages and fields is only recommended. However, any character combination may be replaced by the template engine.

The elements are replaced by the addElement() method into the template class.

We "logically" define 3 type of elements. The separation is only for human logic. The system doesn't make a difference between them. Anything you give to the method will be replaced, however the syntax is. You may define new syntax and new types of elements at will.
The official elements defined for the templates are:

* parameters, which are any piece of information, generally used to build the HTML code.
* language entries, which are any readable information in the language you choose.
* field values, which are generally useful information from a database or data repository.

!!We highly recommend to use the metaelements instead of simple elements.!!


+++ Parameters replacement:

Parameters generally have the following syntax: ~~__PARAMETER__~~
They usually carry pieces of HTML code, for example a color, a size, a tag.

Example:

```
<pre>
&lt;div style="background-color: ~~__BGCOLOR__~~;">
Welcome to my page.<br />
You may use the same parameter as many time you wish.&lt;br />

&lt;span onclick="alert('hello, world');" class="~~__BUTTONCLASS__~~">Click me!&lt;/span>
&lt;span onclick="alert('hello again, world');" class="~~__BUTTONCLASS__~~">Click me again!&lt;/span>

&lt;/div>
</pre>
```


+++ Languages entries

All the languages entries should have the format: ~~##entry##~~
A language entry is generally anything written into your page that does not come from a database, and should adapt to the language of the client visiting your site.
Using the languages entries may depend on the internationalization of your page.
If your page is going to be in a single language forever, you really dont need to use languages entries.
The languages entries generally carry titles, menu options, tables headers etc.

Example:

```
<pre>
&lt;div style="background-color: blue;">
##welcome##&lt;br />
You may use the same parameter as many time you wish.&lt;br />

&lt;span onclick="alert('##hello##');" class="button">##clickme##!&lt;/span>
&lt;span onclick="alert('##helloagain##');" class="button">##clickme## ##again##!&lt;/span>

&lt;/div>
</pre>
```


+++ Field values

Fields values should have the format: {fieldname}
Your fields source can be a database or any other preferred repository data source.
Is it highly recommended to use this syntax for any data field you need to replace in your template.

Example:

```
<pre>
&lt;div style="background-color: blue;">
Welcome to my site&lt;br />
You may use the same parameter as many time you wish.&lt;br />

Today's temperature is {degres} celsius&lt;br />
or {fdegres} farenheit&lt;br />

Is {degres} degres too cold ? Buy a pullover!&lt;br />

&lt;/div>
</pre>
```


++ MetaElements

The metaelements are the recommended way to use with the templates.
They consist into an injection of an associative array of values, called the **data array**, into the template.
The macrolanguage is directly applied on the structure of the data array.

The data array is a nested set of variables and values with the structure you want (there is no construction rules).

You can inject nearly anything into a template metaelements.

Example of a data array to inject:

```
<pre>$array = array(
  'detail' => array(
     'key1' => array('name' => 'Juan', 'status' => 1),
     'key2' => array('name' => 'José', 'status' => 2),
     'key3' => array('name' => 'Pedro', 'status' => 3),
     'key4' => array('name' => 'Phil', 'status' => 1),
     'key5' => array('name' => 'Patrick', 'status' => 2),
  ),
  'param1' => 'blue',
  'param2' => 'red',
  'param3' => '45px',
  'param4' => '100%',
);
</pre>
```

- The **data array** can be any traversable, iterable, countable object too, as of version 1.01.11 and superior.

- You can access directly any data into the array with its relative path (relative to the level you are when the metaelements are applied, see below).

- There are 5 metaelements in the DomCore templates to use the **data array**:
Data, Reference, Loops, Condition and Debug.

The structure of the metaelements in the template must follow the structure of the data array.

+++ ID access: id

The id is 'a-z', 'A-Z', '0-9' and special chars '.-_'
If you use any other character for the id, the compiler will not recognize the keyword as an id and will surely generate errors.

+++ Scope:

When you use an id to point a value, the template will first search into the available ids of the local level.
If no id is found, the it will search into the upper levers if any

Example:
```
<pre>$array = array(
  'detail' => array(
    'data1' => array(
      'data2' => array(
        'key1' => array('appname' => 'Nested App', 'name' => 'Juan', 'status' => 1),
        'key2' => array('name' => 'José', 'status' => 2),
        'appname' => 'DomCore'
      )
    )
  )
);
</pre>

At the level of 'data2', using ~~{{appname}}~~ will get back 'DomCore'
At the level of 'key1', using ~~{{appname}}~~ will get back 'Nested App'
At the level of 'key2', using ~~{{appname}}~~ will get back 'DomCore'
At the level of 'data1', using ~~{{appname}}~~ will get back an empty string


```


+++ Path access: id>id>id>id

At any level into the data array, you can access any entry into the subset array.

For instance, if you have the following data array:

```
<pre>$array = array(
  'detail' => array(
    'data1' => array(
      'data2' => array(
        'key1' => array('name' => 'Juan', 'status' => 1),
        'key2' => array('name' => 'José', 'status' => 2)
      )
    )
  )
);
</pre>
```

Let's suppose we are into a nested metaelements at the 'data1' level. You may want to access directly the 'Juan' entry.
The path will be:

**data2>key1>name**

The José's status value from the root will be:

**detail>data1>data2>key2>status**



+++ Data: ~~{{id}}~~

The data entries are accesible through a macrolanguage keyword: ~~{{id}}~~

The template can work with a strict mode (elements are accesible **only** with the ~~{{..}}~~ syntax).
or in relax mode (any parameter name into the data array will be replaced, !!which is dangerous since you tend to use standard word to create variable names!!).

!!By default as of version 1.01.11, the default mode is strict!!

The **id** can be a direct name, or a path.

Example:
The two last nested metaelements at the 'data1' level will be:

~~{{data2>key1>name}}~~ to access the 'Juan' entry from the 'data1' level.

The José's status value from the root will be:

~~{{detail>data1>data2>key2>status}}~~



+++ References: &&<var2>id</var2>&& and &&<var2>id</var2>:<var1>templateid</var1>&&

Makes a call to a sub template and replace the &&...&& with the result.

If you use &&<var2>id</var2>&&, this is equivalent to &&<var2>id</var2>:<var2>id</var2>&&

The <var1>templateid</var1> is the id of the ~~[[~~<var1>templateid</var1>~~]]~~ to use.
The <var2>id</var2> is the variable id in the valors vector to inject in the template.
If the id exists in the valors vector, then its value is used to replace elements into the subtemplate.
If the id does not exists, then the subtemplate will be resolved only with the main template elements.

The <var2>id</var2> can be a direct name, or a path to access a data into the **data array**.


Example:

```
**Our vector of values:**
<pre>$array = array(
  '<var2>image</var2>' => array('<var5>src</var5>' => '/pics/logo.gif', '<var6>title</var6>' => 'Title of my image')
);
</pre>

**The template: (strict mode)**
<pre>
&&<var1>header</var1>&&
&&<var2>image</var2>:<var3>body</var3>&&
&&<var4>footer</var4>&&

using ~~{{src}}~~ and ~~{{title}}~~ out of the body template is useless, since they are into the ~~{{image}}~~ vector, thus to be used into the 'body' template.

~~[[~~<var1>header</var1>~~]]~~
  Data header&lt;hr />
~~[[]]~~

~~[[~~<var3>body</var3>~~]]~~
  &lt;img src="~~{{~~<var5>src</var5>~~}}~~" title="~~{{~~<var6>title</var6>~~}}~~" />
~~[[]]~~

~~[[~~<var4>footer</var4>~~]]~~
  &lt;hr />Data footer
~~[[]]~~
</pre>
```


+++ Loops: @@<var1>entry</var1>@@ @@<var1>entry</var1>:<var2>template</var2>@@ and @@<var1>entry</var1>:<var2>template</var2>:<var3>check</var3>@@

Makes a call to a subtemplate for each value in the loop vector (like the values of a table).

If you use @@<var1>entry</var1>@@, this is equivalent to @@<var1>entry</var1>:<var2>entry</var2>:@@
If you use @@<var1>entry</var1>:<var2>templateid</var2>@@, this is equivalent to @@<var1>entry</var1>:<var2>templateid</var2>:@@

If '<var1>entry</var1>' does not exists in the values vector, or is empty, or is not a vector, the <var2>templateid</var2> with suffix '.none' will be searched.
If this template does not exists, nothing will be shown.

If '<var1>entry</var1>' is a vector the following templates will be searched for each line, in that order:
- <var2>templateid</var2>.key.[value]  value is the key of the vector line
- <var2>templateid</var2>.sel.[value]  value is the value of the <var3>check</var3> field if it is defined and existing in the vector line
- <var2>templateid</var2>.first if it is the first element of the array set (new from v1.01.11)
- <var2>templateid</var2>.loopalt if the line number is even
- <var2>templateid</var2>.loop
- <var2>templateid</var2>

The <var1>entry</var1> and <var3>check</var3> can be a direct name, or a path to access a data into the data array.


Example:

```
**Our vector of values:**
<pre>$array = array(
  'detail' => array(
     'key1' => array('name' => 'Juan', 'status' => 1),
     'key2' => array('name' => 'José', 'status' => 2),
     'key3' => array('name' => 'Pedro', 'status' => 1),
     'key4' => array('name' => 'Phil', 'status' => 1),
     'key5' => array('name' => 'Patrick', 'status' => 2),
  )
);
</pre>
**The template:**
<pre>
Here comes the loop:<br />
@@detail:eachname:status@@

~~[[eachname.none]]~~
There is nobody in the list<br />
~~[[]]~~

~~%-- First element --%~~
~~[[eachname.first]]~~
&lt;font color="red">Default user:&lt;br />Name: ~~{{name}}~~, Status: ~~{{status}}~~&lt;/font>&lt;br />
~~[[]]~~

~~%-- Number 2 is special --%~~
~~[[eachname.sel.key2]]~~
&lt;font color="red">Name: ~~{{name}}~~, Status: ~~{{status}}~~&lt;/font>&lt;br />
~~[[]]~~

~~[[eachname]]~~
Name: ~~{{name}}~~, Status: ~~{{status}}~~&lt;br />
~~[[]]~~


Another way to show the data:<br />
@@detail@@

~~[[detail.none]]~~
There is nobody in the list&lt;br />
~~[[]]~~

~~[[detail.loopalt]]~~
  &lt;div style="background-color: gray;">Name: ~~{{name}}~~, Status: ??~~{{status}}~~:status??&lt;/div>
  ~~[[status.2]]~~&lt;font color="red">Fired&lt;/font>~~[[]]~~
  ~~[[status]]Ok[[]]~~
~~[[]]~~

~~[[detail.loop]]~~
  &lt;div style="background-color: white;">Name: ~~{{name}}~~, Status: ??~~{{status}}~~:status??&lt;/div>
  ~~[[status.2]]~~&lt;font color="red">Fired&lt;/font>~~[[]]~~
  ~~[[status]]Ok[[]]~~
~~[[]]~~
</pre>
```


+++ Conditional: ??<var1>entry</var1>?? ??<var1>entry</var1>:<var2>templateid</var2>?? and ??<var1>entry</var1>:<var2>templateid</var2>:<var3>check</var3>??

Makes a call to a subtemplate only if the field exists and have a value.

If you use ??<var1>entry</var1>??, this is equivalent to ??<var1>entry</var1>:<var2>entry</var2>:??
If you use ??<var1>entry</var1>:<var2>templateid</var2>??, this is equivalent to ??<var1>entry</var1>:<var2>templateid</var2>:??

If '<var1>entry</var1>' does not exists in the values vector, or is empty, or is not a vector, the <var2>templateid</var2> with suffix '.none' will be searched.
If this template does not exists, nothing will be shown.

The template with the suffix .none is **mandatory**.

If '<var1>entry</var1>' is a vector the following templates will be searched, in that order:
- <var2>templateid</var2>.[value]  value is the value of the <var3>check</var3> field if it is defined and existing in the vector line
- <var2>templateid</var2>

The <var1>entry</var1> and <var3>check</var3> can be a direct name, or a path to access a data into the data array.

Example:

```
**Our vector of values:**
<pre>$array = array(
  'image1' => null,
  'image2' => array('src' => '/pics/logo.gif', 'title' => 'Title of the image')
  'image3' => array('src' => '/pics/logo.gif', 'title' => 'Another title of the image', 'status' => 1)
);
</pre>
**The template:**

??image1:image:status??
??image2:image:status??
??image3:image:status??

~~[[image.none]][[]]~~

~~[[image.1]]~~
Image with status=1:<br />
~~&lt;img src="{{src}}" style="border: 1px solid black;" alt="{{title}}" title="{{title}}" />&lt;br />~~
~~[[]]~~

~~[[image]]~~
~~&lt;img src="{{src}}" alt="{{title}}" title="{{title}}" />&lt;br />~~
~~[[]]~~

```




Exammples:

```
* Load the template file:
<pre>
$buffer = file_get_contents('path/to/your/file.template');
</pre>

* Create the template object with your template file string:
<pre>
$template = new \core\WATemplate($buffer);
</pre>

* Inject elements and metaelements in the template object:
<pre>
$template->addElement('variable', 'value');
$template->addElements(array('variable' => 'value'));
$template->metaElements(array('variable' => 'value'));
</pre>

* Resolve the template to get the generated code:
<pre>
print $template->resolve();
~~//~~ similar to
print $template;
</pre>
```

If you want to use caches and compiled files for much faster access (for instance if you use it in a CMS or so), it is better to use TemplateSource since it resolve all the caches workflow, up to stock the template in shared memory.

How to use it:

```
* Create the template source:
<pre>
$SHM = new \core\WASHM(); ~~//~~ Do no forget to use a unique ID for your application
$templatesource = new \datasources\TemplateSource(
  new \datasources\FileSource('base/path/', 'local/path/', 'your_template_file.template'),
  new \datasources\FastObjectSource(
    new \datasources\FileSource('base/path/', 'local/path/', 'your_afo_file.afo'),
    new \datasources\SHMSource('unique_memory_id', $SHM)
    )
 );
</pre>

* Use the template source to retrieve the template object:
<pre>
$template = $templatesource->read();
</pre>

* Inject elements and metaelements in the template object:
<pre>
$template->addElement('variable', 'value');
$template->addElements(array('variable' => 'value'));
~~//~~ Same as
$template->metaElements(array('variable' => 'value'));
</pre>

* Resolve the template to get the generated code:
<pre>
print $template->resolve();
~~//~~ Similar to
print $template;
</pre>
```

As a reference, using the simple \core\WATemplate object will take approx 12 millisecond to load/compile/resolve the template.
Using the shared memory cache will take only 2 milliseconds to get the template and resolve it (on a 2GHz Xeon processor).

Talking about a good CMS or an application with many templates, using the \datasources\TemplateSource decreases dramatically file accesses and calculation time of your code.











+++ Debug tools

There are two keywords to dump the content of the vector of values, i.e. the elements and the metaelements.
This is very useful when you dont know the code that calls the template, don't remember some values, or for debug facilities.

++++ ~~!!dump!!~~
Shows the totality of the elements and metaelements, variables and values.

++++ ~~!!list!!~~
Shows only the variables names of the elements and metaelements, values are not shown.



3. Functions Reference
------------------------

To use the package:

import "github.com/webability-go/xcore"

type XTemplateParam struct {
  paramtype int
  data string
  children *XTemplateData
}

type XTemplateData []XTemplateParam

type XTemplate struct {
  Name string
  Root *XTemplateData
  SubTemplates map[string]*XTemplate
}

func NewXTemplate() *XTemplate {

func NewXTemplateFromFile(file string) (*XTemplate, error) {

func NewXTemplateFromString(data string) (*XTemplate, error) {

func (t *XTemplate)LoadFile(file string) error {

func (t *XTemplate)LoadString(data string) error {

func (t *XTemplate)compile(data string) error {

func (t *XTemplate)AddTemplate(name string, tmpl *XTemplate) {

func (t *XTemplate)GetTemplate(name string) *XTemplate {

func (t *XTemplate)Execute(data XDatasetDef) string {

func (t *XTemplate)injector ( datacol XDatasetCollectionDef, language *XLanguage ) string {

func searchConditionValue(id string, data XDatasetCollectionDef) string {

func buildValue(data interface{}) string {

func (t *XTemplate)Print() string {


---
*/
//
package xcore

// VERSION is the used version nombre of the XCore library.
const VERSION = "0.3.1"

// LOG is the flag to activate logging on the library.
// if LOG is set to TRUE, LOG indicates to the XCore libraries to log a trace of functions called, with most important parameters.
// LOG can be set to true or false dynamically to trace only parts of code on demand.
var LOG = false
