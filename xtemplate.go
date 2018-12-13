package xcore

import (
  "fmt"
  "os"
  "io/ioutil"
  "strings"
  "regexp"
  "errors"
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
  MetaReference              = 3        // an URL variable coming through a query ?variable=value
  MetaRange                  = 4        // Parameter passed to the page Run by code
  MetaCondition              = 5        // System (site) parameter
  MetaDump                   = 6        // Main page called parameters (into .page file)
  MetaVariable               = 7        // this page parameters (into .page file), same as Main page parameters if it's the external called page

  MetaTemplateStart          = 101      // Temporal nested box start tag
  MetaTemplateEnd            = 102      // Temporal nested box end tag
  
  MetaUnused                 = -1       // a "not used anymore" param to be freed
)

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
  return &XTemplate{}
}

func NewXTemplateFromFile(file string) (*XTemplate, error) {
  t := &XTemplate{}
  err := t.LoadFile(file)
  if err != nil { return nil, err }
  return t, nil
}

func NewXTemplateFromString(data string) (*XTemplate, error) {
  t := &XTemplate{}
  err := t.LoadString(data)
  if err != nil { return nil, err }
  return t, nil
}

func (t *XTemplate)LoadFile(file string) error {
  tFile, err := os.Open(file)
  if err != nil { return err }
  data, err := ioutil.ReadAll(tFile)
  if err != nil { return err }
  err = tFile.Close()
  if err != nil { return err }
  return t.LoadString(string(data))
}

func (t *XTemplate)LoadString(data string) error {
  return t.compile(data)
}

func (t *XTemplate)compile(data string) error {
  // build, compile return result
  code :=
      `(?s)`+                                                                             // . is multiline
      
      // ==== COMENTS
      `(%)--(.*?)--%\n?`+                                                                 // index based 1

      // ==== LANGUAGE INJECTION
      `|(#)#(.+?)##`+                                                                     // index based 3

      // ==== ELEMENTS
      `|(&)&(.+?)&&`+                                                                     // index based 5
      `|(@)@(.+?)@@`+                                                                     // index based 7
      `|(\?)\?(.+?)\?\?`+                                                                 // index based 9
      `|(\!)\!(.+?)\!\!`+                                                                 // index based 11
      `|(\{)\{(.+?)\}\}`+                                                                 // index based 13

      // ==== NESTED ELEMENTS (SUB TEMPLATES)
      `|\[\[(\])\]`+                                                                      // index based 15
      `|(\[)\[(.+?)\]\]`                                                                  // index based 16

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
      
      // BUILD THE 2 & PARTS
      param.data = matches[i][6]


    } else if matches[i][7] == "@" {
      param.paramtype = MetaRange             // pageparam
      
      // BUILD THE 3 @ PARTS
      param.data = matches[i][8]
      
      
    } else if matches[i][9] == "?" {
      param.paramtype = MetaCondition         // local pageparam

      // BUILD THE 3 ? PARTS
      param.data = matches[i][10]
      
      
    } else if matches[i][11] == "!" {
      param.paramtype = MetaDump              // instance param
      param.data = matches[i][12]
    } else if matches[i][13] == "{" {
      param.paramtype = MetaVariable          // local instance param
      param.data = matches[i][14]
    } else if matches[i][16] == "[" {
      param.paramtype = MetaTemplateStart     // javascript call for header
      param.data = matches[i][17]
    } else if matches[i][15] == "]" {
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

  // second pass: all the sub templates into the Subtemplates
  startpointers := []int{}
  subtemplates := []*XTemplate{}
  actualtemplate := t
  for i, x := range compiled {
    if x.paramtype == MetaTemplateStart {
      startpointers = append(startpointers, i)
      subtemplates = append(subtemplates, actualtemplate)
      actualtemplate = &XTemplate{Name: x.data, Root: nil,}
    } else if x.paramtype == MetaTemplateEnd {
      // we found the end of the nested box, lets create a nested param array from stacked startpointer up to i
      last := len(startpointers)-1
      if last < 0 {
        return errors.New("Error: template mismatch Start/End")
      }
      startpointer := startpointers[last]
      startpointers = startpointers[:last]
      
      var subset XTemplateData
      for ptr := startpointer+1; ptr < i; ptr++ {   // we ignore the BOX]] end param (we dont need it in the hierarchic structure)
        if compiled[ptr].paramtype != MetaUnused { // we just ignore params marked to be deleted
          subset = append(subset, compiled[ptr])
          compiled[ptr].paramtype = MetaUnused   // marked to be deleted, traslated to a substructure
        }
      }
      actualtemplate.Root = &subset

      uppertemplate := subtemplates[last]
      subtemplates = subtemplates[:last]
      
      // If there are |, we separate the templates and add every one to the list with same pointer
      pospipe := strings.Index(actualtemplate.Name, "|")
      if pospipe >= 0 {
        vals := strings.Split(actualtemplate.Name, "|")
        for _, v := range vals {
          if len(v) > 0 {
            uppertemplate.AddTemplate(v, actualtemplate)
          }
        }
      } else {
        uppertemplate.AddTemplate(actualtemplate.Name, actualtemplate)
      }
      
      // pop actualtemplate
      actualtemplate = uppertemplate;

      compiled[startpointer].paramtype = MetaUnused  // marked to be deleted, no need of start template
      compiled[i].paramtype = MetaUnused   // marked to be deleted, no need of end template
    }
  }
  if len(startpointers) > 0 {
    return errors.New("Error: template mismatch Start/End")
  }
  
  // last pass: delete params marked to be deleted and concatenate strings
  currentpointer := 0
  for i, x := range compiled {
    if x.paramtype != MetaUnused {
      if currentpointer != i {
        compiled[currentpointer] = x
      }
      currentpointer += 1
    }
  }
  
  // ************** NOTE: MISSING OPEN/CLOSE SUBTEMPLATES = COMPILATION ERROR
  
  compiled = compiled[:currentpointer]
  t.Root = &compiled
  return nil
}

func (t *XTemplate)AddTemplate(name string, tmpl *XTemplate) {
  if t.SubTemplates == nil {
    t.SubTemplates = make(map[string]*XTemplate)
  }
  t.SubTemplates[name] = tmpl
}

func (t *XTemplate)GetTemplate(name string) *XTemplate {
  if t.SubTemplates == nil { return nil }
  return t.SubTemplates[name]
}

func (t *XTemplate)Execute(data XDataset) string {
  // Does data has a language ?
  lang, ok := data["#"]
  var language *XLanguage = nil
  if ok { language = lang.(*XLanguage) }
  stack := []XDataset{data}
  return t.injector(stack, language)
}

func (t *XTemplate)injector ( data []XDataset, language *XLanguage ) string {
  var injected []string
  for _, v := range *t.Root {
    switch v.paramtype {
      case MetaString: // included string from original code
        injected = append(injected, v.data)
      case MetaComment:
        // nothing to do: comment ignored
      case MetaLanguage:
        if language != nil {
          injected = append(injected, language.Get(v.data))
        }
      case MetaReference:     // Reference
        // search for the subtemplate and reinject anything into it
        subt := t.GetTemplate(v.data)
        if subt != nil {
          // if v.data is a substructure into data, then we stack the data and inject new stacked data
          
          substr := subt.injector(data, language)
          injected = append(injected, substr)
        }
      case MetaVariable:
        value := searchValue(v.data, data)
        injected = append(injected, value)
      case MetaRange:    // Range (loop over subset)

        // ***** subtemplate depends on first, last, loop # counter, key, condition etc
        subt := t.GetTemplate(v.data)

        if subt != nil {
          value := searchRangeValue(v.data, data)
          for _, ds := range value {
            // if v.data is a substructure into data, then we stack the data and inject new stacked data
            data = append(data, ds)
            substr := subt.injector(data, language)
            injected = append(injected, substr)
            // unstack ds
            data = data[:len(data)-1]
          }
        }
      case MetaCondition:

        // ***** subtemplate depends on condition completion
        subt := t.GetTemplate(v.data)
        value := searchConditionValue(v.data, data)

        // if v.data is a substructure into data, then we stack the data and inject new stacked data
        if value != "" {
          substr := subt.injector(data, language)
          injected = append(injected, substr)
        }
      case MetaDump:
        substr := dump(data[0])
        injected = append(injected, substr)
      default:
        injected = append(injected, "THE METALANGUAGE FROM OUTERSPACE IS NOT SUPPORTED: " + fmt.Sprint(v.paramtype))
    }
  }
  // return the page string
  return strings.Join(injected, "")
}

func searchRangeValue(id string, data []XDataset) XRangeDataset {
  // scan data for each dataset in order top to bottom
  for i := len(data)-1; i >= 0; i-- {
    val := scanValue(id, data[i])
    if val != nil {
      return val.(XRangeDataset)
    }
  }
  return nil
}

func searchConditionValue(id string, data []XDataset) string {
  // scan data for each dataset in order top to bottom
  for i := len(data)-1; i >= 0; i-- {
    val := scanValue(id, data[i])
    if val != nil {
      return buildValue(val)
    }
  }
  return ""
}

func searchValue(id string, data []XDataset) string {
  // scan data for each dataset in order top to bottom
  for i := len(data)-1; i >= 0; i-- {
    val := scanValue(id, data[i])
    if val != nil {
      return buildValue(val)
    }
  }
  return ""
}

func scanValue(id string, data XDataset) interface{} {
  // scan data for the id
  possup := strings.Index(id, ">")
  if possup >= 0 {
    first := id[:possup]
    // check limits: first == "", second part == ""
    if entry, ok := data[first]; ok {
      // if entry IS map[string]interface{} entonces podemos seguir en la estructura
      // Check also if it's a function that returns a map[string]interface{}
      return scanValue(id[possup+1:], entry.(XDataset))
    }
    return nil
  } else {
    if entry, ok := data[id]; ok {
      return entry
    }
  }
  return nil
}

func buildValue(data interface{}) string {
  // if data is string, return data
  // if data is a type, return conversion
  // if data is a function, call the function then return the value
  return fmt.Sprint(data)
}

func (t *XTemplate) Print() string {
  return fmt.Sprint(t)
}

func dump(data XDataset) string {
  return "DUMP OF DATASET"
}
