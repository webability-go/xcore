package xcore

import (
  "fmt"
  "os"
  "io/ioutil"
  "strings"
  "regexp"
)

/*
class to compile and keep a Template string
A template is a set of HTML/XML (or any other language) set of:

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
   ??xx??   if/then/else
   @@xx@@   loops
   &&xx&&   references
   !!xx!!   debug (dump)
*/

const (
  MetaString                 = 0        // a simple string to integrate into the code
  MetaComment                = 1        // Comment, ignore it
  
  MetaLanguage               = 2        // one param of the URL parameters list, index-1 based [page]/value1/value2...
  MetaReference            = 3        // an URL variable coming through a query ?variable=value
  MetaRange                  = 4        // Parameter passed to the page Run by code
  MetaCondition               = 5        // System (site) parameter
  MetaDump              = 6        // Main page called parameters (into .page file)
  MetaVariable         = 7        // this page parameters (into .page file), same as Main page parameters if it's the external called page
  MetaTemplate                    = 15       // Nested box with inner data

  MetaTemplateStart      = 101      // Temporal nested box start tag
  MetaTemplateEnd        = 102      // Temporal nested box end tag
  
  MetaUnused                 = -1       // a "not used anymore" param to be freed
)

type XTemplate struct {
  Name string
  Root *XTemplateData
  SubTemplates map[string]*XTemplate
}

func NewXTemplate() *XTemplate {
  return &XTemplate{}
}

func (t *XTemplate)LoadFile(file string) error {
  xmlFile, _ := os.Open(file)
  defer xmlFile.Close()
  byteValue, _ := ioutil.ReadAll(xmlFile)
  return t.compile(string(byteValue))
}

func (t *XTemplate)LoadString(data string) error {
  return t.compile(data)
}

type XTemplateParam struct {
  paramtype int
  data string
  children *XTemplateData
}

type XTemplateData []XTemplateParam

func (t *XTemplate)compile(data string) error {
  // build, compile return result
  code :=
      `(?s)`+                                                                             // . is multiline
      
      // ==== COMENTS
      `(%)--(.*?)--%\n?`+                                                                 // index based 1

      // ==== LANGUAGE INJECTION
      `|(#)#(.*?)##`+                                                                     // index based 3

      // ==== ELEMENTS
      `|(&)&(.*?)&&`+                                                                     // index based 5
      `|(@)@(.*?)@@`+                                                                     // index based 7
      `|(\?)\?(.*?)\?\?`+                                                                 // index based 9
      `|(\!)\!(.*?)\!\!`+                                                                 // index based 11
      `|(\{)\{(.*?)\}\}`+                                                                 // index based 13

      // ==== NESTED ELEMENTS (SUB TEMPLATES)
      `|(\[)\[(.*?)\]\]`+                                                                 // index based 15
      `|\[\[(\])\]`                                                                       // index based 17

  codex := regexp.MustCompile(code)
  indexes := codex.FindAllStringIndex(data, -1)
  matches := codex.FindAllStringSubmatch(data, -1)

  var compiled XTemplateData
  pointer := 0
  for i, x := range indexes {
    if pointer != x[0] {
      compiled = append(compiled, *(&XTemplateParam{paramtype: MetaString, data: data[pointer:x[0]],}))
    }

    param := &XTemplateParam{}
    if matches[i][1] == "%" {
      param.paramtype = MetaComment           // comment
      param.data = matches[i][2]
    } else if matches[i][3] == "#" {
      param.paramtype = MetaLanguage          // Entry Param
      param.data = matches[i][4]
    } else if matches[i][5] == "&" {
      param.paramtype = MetaReference         // sysparam
      param.data = matches[i][6]
    } else if matches[i][7] == "@" {
      param.paramtype = MetaRange             // pageparam
      param.data = matches[i][8]
    } else if matches[i][9] == "?" {
      param.paramtype = MetaCondition         // local pageparam
      param.data = matches[i][10]
    } else if matches[i][11] == "!" {
      param.paramtype = MetaDump              // instance param
      param.data = matches[i][12]
    } else if matches[i][13] == "{" {
      param.paramtype = MetaVariable          // local instance param
      param.data = matches[i][14]
    } else if matches[i][15] == "[" {
      param.paramtype = MetaTemplateStart     // javascript call for header
      param.data = matches[i][16]
    } else if matches[i][17] == "]" {
      param.paramtype = MetaTemplateEnd       // css call for header
    } else {
      param.paramtype = MetaUnused            // unknown, will be removed
    }
    compiled = append(compiled, *param)
    pointer = x[1]
  }
  // end of data
  if pointer != len(data) {
    compiled = append(compiled, *(&XTemplateParam{paramtype: MetaString, data: data[pointer:],}))
  }
  
  // second pass: all the nested boxes goes into a subset of subtemplates
  startpointers := []int{}
  for i, x := range compiled {
    if x.paramtype == MetaTemplateStart {
      startpointers = append(startpointers, i)
    } else if x.paramtype == MetaTemplateEnd {
      // we found the end of the nested box, lets create a nested param array from stacked startpointer up to i
      last := len(startpointers)-1
      startpointer := startpointers[last]
      startpointers = startpointers[:last]
      
      var subset XTemplateData
      for ptr := startpointer+1; ptr < i; ptr++ {   // we ignore the BOX]] end param (we dont need it in the hierarchic structure)
        if compiled[ptr].paramtype != MetaUnused { // we just ignore params marked to be deleted
          subset = append(subset, compiled[ptr])
          compiled[ptr].paramtype = MetaUnused   // marked to be deleted, traslated to a substructure
        }
      }
      compiled[startpointer].paramtype = MetaTemplate
      compiled[startpointer].children = &subset
      compiled[i].paramtype = MetaUnused   // marked to be deleted, on need of end box
    }
  }
  
  // last pass: delete params marked to be deleted
  currentpointer := 0
  for i, x := range compiled {
    if x.paramtype != MetaUnused {
      if currentpointer != i {
        compiled[currentpointer] = x
      }
      currentpointer += 1
    }
  }
  compiled = compiled[:currentpointer]
  t.Root = &compiled
  return nil
}

func (t *XTemplate) GetTemplate( name string ) *XTemplate {
  return nil
}

func (t *XTemplate) Execute( data interface{} ) string {
  return t.Root.Execute(data)
}

func (c *XTemplateData) Execute ( data interface{} ) string {
  // third pass: inject meta language
  var injected []string
  for _, v := range *c {
    switch v.paramtype {
      case MetaString: // included string from original code
        injected = append(injected, v.data)
      case MetaComment:
        // nothing to do: comment ignored
      case MetaLanguage:
//        injected = append(injected, "LANGUAGE ENTRY " + v.data)
      case MetaReference:     // Reference
//        injected = append(injected, "JS CALL NOT IMPLEMENTED YET: " + v.data)
      case MetaRange:    // Range (loop over subset)
//        injected = append(injected, "CSS CALL NOT IMPLEMENTED YET: " + v.data)
      case MetaCondition:
        // build the params
      case MetaDump:
        // build the params
      case MetaVariable:
        // build the params
      case MetaTemplate:
        // build the template
      default:
        injected = append(injected, "THE METALANGUAGE FROM OUTERSPACE IS NOT SUPPORTED: " + fmt.Sprint(v.paramtype))
    }
  }
  // return the page string
  return strings.Join(injected, "")
}
