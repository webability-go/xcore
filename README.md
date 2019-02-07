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

V0.0.6 - 2019-
-----------------------
- Added xcache.GetId()

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



Manual:
=======================

XCache
=======================
XCache is a library to cache all the data you want into current application memory for a very fast access to the data.
The access to the data support multithreading and concurrency. For the same reason, this type of cache is not persistant (if you exit the application)
and cannot grow too much (as memory is the limit).
However, you can control a timeout of each cache piece, and eventually the comparison against a source (file, database, etc) to invalid the cache.

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
  myfiles.Set(", "somedata")
  myfiles.Set("/home/sites/file2.txt", "somedata")

  go somefunc()
}

func somefun() {
  data, invalid := myfiles.Get("https://developers.webability.info:82/");
  
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

  
  
XDataSet
=======================

1. Overview
------------------------

The XDataSetDef is an interfase to build a standard set of data optionally nested and hierarchical, that can be used for any purpose:
- Keep complex data in memory
- Create JSON structures
- Inject data into templates
- Interchange database data (records set and record)
etc

You can store into it generic supported data, as well as any complex interface structures:
- Int
- Float
- String
- Time
- Bool
- []Int
- []Float
- []String
- []Time
- []Bool
- XDataSetDef (anything extended with this interface)
- XDataSetCollectionDef (anything extended with this interface)
- Anything else ( interface{} )

Example:



2. Reference
------------------------

type XDatasetDef interface {
  // Stringify will dump the content into a human readable string
  Stringify() string

  // Set will associate the data to the key. If it already exists, it will be replaced
  Set(key string, data interface{})

  // Get will return the value associated to the key if it exists, or bool = false
  Get(key string) (interface{}, bool)
  // Same as Get but will return the value associated to the key as a XDatasetDef if it exists, or bool = false
  GetDataset(key string) (XDatasetDef, bool)
  // Same as Get but will return the value associated to the key as a XDatasetCollectionDef if it exists, or bool = false
  GetCollection(key string) (XDatasetCollectionDef, bool)

  // Same as Get but will return the value associated to the key as a string if it exists, or bool = false
  GetString(key string) (string, bool)
  // Same as Get but will return the value associated to the key as a string if it exists, or bool = false
  GetBool(key string) (bool, bool)
  // Same as Get but will return the value associated to the key as a string if it exists, or bool = false
  GetInt(key string) (int, bool)
  // Same as Get but will return the value associated to the key as a string if it exists, or bool = false
  GetFloat(key string) (float64, bool)
  // Same as Get but will return the value associated to the key as a XDatasetCollectionDef if it exists, or bool = false
  GetStringCollection(key string) ([]string, bool)
  // Same as Get but will return the value associated to the key as a string if it exists, or bool = false
  GetBoolCollection(key string) ([]bool, bool)
  // Same as Get but will return the value associated to the key as a string if it exists, or bool = false
  GetIntCollection(key string) ([]int, bool)
  // Same as Get but will return the value associated to the key as a string if it exists, or bool = false
  GetFloatCollection(key string) ([]float64, bool)
  
  // Del will delete the data associated to the key and deletes the key entry
  Del(key string)
}

/* Basic dataset */
type XDataset map[string]interface{}


type XDatasetCollectionDef interface {
  // Stringify will dump the content into a human readable string
  Stringify() string

  // Will add a datasetdef to the beginning of the collection
  Unshift(data XDatasetDef)
  // Will remove the first datasetdef of the collection and return it
  Shift() XDatasetDef

  // Will add a datasetdef to the end of the collection
  Push(data XDatasetDef)
  // Will remove the last datasetdef of the collection and return it
  Pop() XDatasetDef

  // Will count the quantity of entries
  Count() int

  // Will get the entry by the index and let it in the collection
  Get(index int) (XDatasetDef, bool)

  // Will search for the data associated to the key by priority (last entry is the most important)
  // returns bool = false if nothing has been found
  GetData(key string) (interface{}, bool)

  // Same as GetData but will convert the result to a string if possible
  // returns bool = false if nothing has been found
  GetDataString(key string) (string, bool)
  // Same as GetData but will convert the result to an int if possible
  // returns bool = false if nothing has been found
  GetDataInt(key string) (int, bool)
  // Same as GetData but will convert the result to a boolean if possible
  // returns second bool = false if nothing has been found
  GetDataBool(key string) (bool, bool)
  // Same as GetData but will convert the result to a float if possible
  // returns bool = false if nothing has been found
  GetDataFloat(key string) (float64, bool)
  // Same as GetData but will convert the result to a collection of data if possible
  // returns bool = false if nothing has been found
  GetCollection(key string) (XDatasetCollectionDef, bool)
}

// =====================
// XDatasetConnection
// =====================

type XDatasetCollection []XDatasetDef




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
A template is a set of HTML/XML (or any other language) string with a meta language to inject variables and build a final string.

The XCore XTemplate system is based on the injection of parameters, language translation strings and data fields directly into the HTML (Or any other language you need) template.
The HTML itself (or any other language) is a fixed code not directly used by the template system, but used to dress the data you want to represent in your prefered language.

The variables to inject are into a XDataSet structure or into a structure extended from XDataSetDef interfase.
The injection of data is based on a XDataSet structure of values that can be nested into another XDataSet and XDataSetConnection and so on.
The template compiler recognize nested arrays to automatically make loops on the information.

The macrolanguage is extremely simple and is made to be usefull and **really** separate programation from template code (not like other may generic template systems that just mix code and data).

!!Templates are made to store reusable HTML code, and overall easily changeable by **NON PROGRAMMING PEOPLE**.!!

In sight to create and use templates, you have all those possible options to use:

* Comments
* Nested templates, to store many pieces of HTML
* Simple elements, to replace by values in the template. There are various types of simple elements:
** Parameters
** Language entries
** Fields of values
* Meta elements, to build a code based on a data array. There are various types of meta elements:
** Data access with ~~{{...}}~~, to show the value of a data into the data array
** Subtemplates access with ~~&&...&&~~, to call a subtemplate based on the value of an entry in the data array.
** Conditional access with ~~??...??~~, to show a piece of HTML based on the existance or value of an entry in the data array.
** Loops access with ~~@@...@@~~, to repeat a piece of HTML based on a table of values.
** Debug tools with ~~!!...!!~~, to show the data array.


Exammples:

[1]
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
[_]

If you want to use caches and compiled files for much faster access (for instance if you use it in a CMS or so), it is better to use TemplateSource since it resolve all the caches workflow, up to stock the template in shared memory.

How to use it:

[1]
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
[_]

As a reference, using the simple \core\WATemplate object will take approx 12 milisecond to load/compile/resolve the template.
Using the shared memory cache will take only 2 milliseconds to get the template and resolve it (on a 2GHz Xeon processor).

Talking about a good CMS or an application with many templates, using the \datasources\TemplateSource decreases dramatically file accesses and calculation time of your code.





2. Meta Language Reference
------------------------

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
   
++ Comments

You may use comments into your template.
The comments will be discarded immediatly at the compilation of the template and not interfere with the rest of your code.

Comments are defined by ~~%--~~ and ~~--%~~

Example:

[1]<pre>

~~%--~~ This is a comment. It will not appear in the final code. ~~--%~~

~~%--~~
This subtemplate will not be compiled, usable or even visible since it is into a comment
~~[[templateid]]~~
Anything here
~~[[]]~~
~~--%~~

</pre>[_]


++ Nested Templates

You can define new nested templates into your main template
A nested template is defined by:

[1]<pre>
~~[[templateid]]~~
your nested template
~~[[]]~~
</pre>[_]

The id is any combination of letters (a-z, A-Z, accents are welcome too), numbers (0-9), and 3 special chars: .-_

The old syntax also work and is deprecated. It will be definitively removed as of beginning of 2013.

[1]<pre>
%%SUBTEMPLATE(templateid)%%
your nested template
%%ENDSUBTEMPLATE%%
</pre>[_]

There is no limits into nesting templates.
Any nested template will inheritate all the father elements and can use father elements too.

Example:

[1]
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
[_]

You may use more than one id into the same template to avoid repetition of the same code.
The different id's are separated with a pipe |

[1]<pre>
~~[[templateid|anotherid|something.key|andmoreid]]~~
your nested template
~~[[]]~~
</pre>[_]




++ Elements

The DomCore template system is based on a parameters replacement and a simple macrolanguage keywords.
Note: The syntax of the parameters, languages and fields is only recommended. However, any character combination may be replaced by the template engine.

The elements are replaced by the addElement() method into the template class.

We "logically" define 3 type of elements. The separation is only for human logic. The system doesn't make a difference between them. Anything you give to the method will be replaced, however the syntax is. You may define new syntax and new types of elements at will.
The official elements defined for the templates are:

* parameters, which are any piece of information, generaly used to build the HTML code.
* language entries, which are any readable information in the language you choose.
* field values, which are generally usefull information from a database or data repository.

!!We highly recommend to use the metaelements instead of simple elements.!!


+++ Parameters replacement:

Parameters generally have the following syntax: ~~__PARAMETER__~~ 
They usually carry pieces of HTML code, for example a color, a size, a tag.

Example:

[1]<pre>
&lt;div style="background-color: ~~__BGCOLOR__~~;">
Welcome to my page.<br />
You may use the same parameter as many time you wish.&lt;br />

&lt;span onclick="alert('hello, world');" class="~~__BUTTONCLASS__~~">Click me!&lt;/span>
&lt;span onclick="alert('hello again, world');" class="~~__BUTTONCLASS__~~">Click me again!&lt;/span>

&lt;/div>
</pre>[_]


+++ Languages entries

All the languages entries should have the format: ~~##entry##~~
A language entry is generally anything written into your page that does not come from a database, and should adapt to the language of the client visiting your site.
Using the languages entries may depend on the internationalization of your page.
If your page is going to be in a single language forever, you really dont need to use languages entries.
The languages entries generally carry titles, menu options, tables headers etc.

Example:

[1]<pre>
&lt;div style="background-color: blue;">
##welcome##&lt;br />
You may use the same parameter as many time you wish.&lt;br />

&lt;span onclick="alert('##hello##');" class="button">##clickme##!&lt;/span>
&lt;span onclick="alert('##helloagain##');" class="button">##clickme## ##again##!&lt;/span>

&lt;/div>
</pre>[_]


+++ Field values

Fields values should have the format: {fieldname}
Your fields source can be a database or any other prefered repository data source.
Is it highly recomended to use this syntax for any data field you need to replace in your template.

Example:

[1]<pre>
&lt;div style="background-color: blue;">
Welcome to my site&lt;br />
You may use the same parameter as many time you wish.&lt;br />

Today's temperature is {degres} celsius&lt;br />
or {fdegres} farenheit&lt;br />

Is {degres} degres too cold ? Buy a pullover!&lt;br />

&lt;/div>
</pre>[_]


++ MetaElements

The metaelements are the recommended way to use with the templates.
They consist into an injection of an associative array of values, called the **data array**, into the template.
The macrolanguage is directly applied on the structure of the data array.

The data array is a nested set of variables and values with the structure you want (there is no construction rules).

You can inject nearly anything into a template metaelements.

Example of a data array to inject:

[1]
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
[_]

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
[1]
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


[_]
 

+++ Path access: id>id>id>id

At any level into the data array, you can access any entry into the subset array.

For instance, if you have the following data array:

[1]
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
[_]

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

[2]
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
[_]


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

[2]
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
[_]


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

[2]
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

[_]


+++ Debug tools

There are two keywords to dump the content of the vector of values, i.e. the elements and the metaelements.
This is very usefull when you dont know the code that calls the template, don't remember some values, or for debug facilities.

++++ ~~!!dump!!~~
Shows the totality of the elements and metaelements, variables and values.

++++ ~~!!list!!~~
Shows only the variables names of the elements and metaelements, values are not shown.



3. Functions Reference
------------------------

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
