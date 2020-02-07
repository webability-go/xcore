package xcore

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
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

// MetaString and other consts:
//   type of elements present in the template
const (
	MetaString  = 0 // a simple string to integrate into the code
	MetaComment = 1 // Comment, ignore it

	MetaLanguage  = 2 // one param of the URL parameters list, index-1 based [page]/value1/value2...
	MetaReference = 3 // an URL variable coming through a query ?variable=value
	MetaRange     = 4 // Parameter passed to the page Run by code
	MetaCondition = 5 // System (site) parameter
	MetaDump      = 6 // Main page called parameters (into .page file)
	MetaVariable  = 7 // this page parameters (into .page file), same as Main page parameters if it's the external called page

	MetaTemplateStart = 101 // Temporal nested box start tag
	MetaTemplateEnd   = 102 // Temporal nested box end tag

	MetaUnused = -1 // a "not used anymore" param to be freed
)

// XTemplateParam is a parameter definition into the template
type XTemplateParam struct {
	ParamType int
	Data      string
	//	children  *XTemplateData
}

// XTemplateData is an Array of all the parameters into the template
type XTemplateData []XTemplateParam

// XTemplate is the plain template structure
type XTemplate struct {
	Name         string
	Root         *XTemplateData
	SubTemplates map[string]*XTemplate
}

// NewXTemplate will create a new empty template
func NewXTemplate() *XTemplate {
	return &XTemplate{}
}

// NewXTemplateFromFile will create a new template from a file containing the template code
func NewXTemplateFromFile(file string) (*XTemplate, error) {
	t := &XTemplate{}
	err := t.LoadFile(file)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// NewXTemplateFromString will create a new template from a string containing the template code
func NewXTemplateFromString(data string) (*XTemplate, error) {
	t := &XTemplate{}
	err := t.LoadString(data)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// LoadFile will load a file into the template
func (t *XTemplate) LoadFile(file string) error {
	tFile, err := os.Open(file)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(tFile)
	if err != nil {
		return err
	}
	err = tFile.Close()
	if err != nil {
		return err
	}
	return t.LoadString(string(data))
}

// LoadString will load a string into the template
func (t *XTemplate) LoadString(data string) error {
	return t.compile(data)
}

// compile will interprete the template code into objects
func (t *XTemplate) compile(data string) error {
	// build, compile return result
	code :=
		`(?s)` + // . is multiline

			// ==== COMENTS
			`(%)--(.*?)--%\n?` + // index based 1

			// ==== LANGUAGE INJECTION
			`|(#)#(.+?)##` + // index based 3

			// ==== ELEMENTS
			`|(&)&(.+?)&&` + // index based 5
			`|(@)@(.+?)@@` + // index based 7
			`|(\?)\?(.+?)\?\?` + // index based 9
			`|(\!)\!(.+?)\!\!` + // index based 11
			`|(\{)\{(.+?)\}\}` + // index based 13

			// ==== NESTED ELEMENTS (SUB TEMPLATES)
			`|\[\[(\])\]` + // index based 15
			`|(\[)\[([a-z0-9\.\-_]+?)\]\]` // index based 16

	codex := regexp.MustCompile(code)
	indexes := codex.FindAllStringIndex(data, -1)
	matches := codex.FindAllStringSubmatch(data, -1)

	var compiled XTemplateData
	pointer := 0
	for i, x := range indexes {
		if pointer != x[0] {
			compiled = append(compiled, *(&XTemplateParam{ParamType: MetaString, Data: data[pointer:x[0]]}))
		}

		param := &XTemplateParam{}
		if matches[i][1] == "%" {
			param.ParamType = MetaComment // comment
			param.Data = matches[i][2]
		} else if matches[i][3] == "#" {
			param.ParamType = MetaLanguage // Entry Param
			param.Data = matches[i][4]
		} else if matches[i][5] == "&" {
			param.ParamType = MetaReference // sysparam

			// BUILD THE 2 & PARTS
			param.Data = matches[i][6]

		} else if matches[i][7] == "@" {
			param.ParamType = MetaRange // pageparam

			// BUILD THE 3 @ PARTS
			param.Data = matches[i][8]

		} else if matches[i][9] == "?" {
			param.ParamType = MetaCondition // local pageparam

			// BUILD THE 3 ? PARTS
			param.Data = matches[i][10]

		} else if matches[i][11] == "!" {
			param.ParamType = MetaDump // instance param
			param.Data = matches[i][12]
		} else if matches[i][13] == "{" {
			param.ParamType = MetaVariable // local instance param
			param.Data = matches[i][14]
		} else if matches[i][16] == "[" {
			param.ParamType = MetaTemplateStart // javascript call for header
			param.Data = matches[i][17]
		} else if matches[i][15] == "]" {
			param.ParamType = MetaTemplateEnd // css call for header
		} else {
			param.ParamType = MetaUnused // unknown, will be removed
		}
		compiled = append(compiled, *param)
		pointer = x[1]
	}
	// end of Data
	if pointer != len(data) {
		compiled = append(compiled, *(&XTemplateParam{ParamType: MetaString, Data: data[pointer:]}))
	}

	// second pass: all the sub templates into the Subtemplates
	startpointers := []int{}
	subtemplates := []*XTemplate{}
	actualtemplate := t
	for i, x := range compiled {
		if x.ParamType == MetaTemplateStart {
			startpointers = append(startpointers, i)
			subtemplates = append(subtemplates, actualtemplate)
			actualtemplate = &XTemplate{Name: x.Data, Root: nil}
		} else if x.ParamType == MetaTemplateEnd {
			// we found the end of the nested box, lets create a nested param array from stacked startpointer up to i
			last := len(startpointers) - 1
			if last < 0 {
				return errors.New("Error: template mismatch Start/End")
			}
			startpointer := startpointers[last]
			startpointers = startpointers[:last]

			var subset XTemplateData
			for ptr := startpointer + 1; ptr < i; ptr++ { // we ignore the BOX]] end param (we dont need it in the hierarchic structure)
				if compiled[ptr].ParamType != MetaUnused { // we just ignore params marked to be deleted
					subset = append(subset, compiled[ptr])
					compiled[ptr].ParamType = MetaUnused // marked to be deleted, traslated to a substructure
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
			actualtemplate = uppertemplate

			compiled[startpointer].ParamType = MetaUnused // marked to be deleted, no need of start template
			compiled[i].ParamType = MetaUnused            // marked to be deleted, no need of end template
		}
	}
	if len(startpointers) > 0 {
		return errors.New("Error: template mismatch Start/End")
	}

	// last pass: delete params marked to be deleted and concatenate strings
	currentpointer := 0
	for i, x := range compiled {
		if x.ParamType != MetaUnused {
			if currentpointer != i {
				compiled[currentpointer] = x
			}
			currentpointer++
		}
	}

	// TODO(phil) NOTE: MISSING OPEN/CLOSE SUBTEMPLATES = COMPILATION ERROR

	compiled = compiled[:currentpointer]
	t.Root = &compiled
	return nil
}

// AddTemplate will add a sub template to this template
func (t *XTemplate) AddTemplate(name string, tmpl *XTemplate) {
	if t.SubTemplates == nil {
		t.SubTemplates = make(map[string]*XTemplate)
	}
	t.SubTemplates[name] = tmpl
}

// GetTemplate gets a sub template existing into this template
func (t *XTemplate) GetTemplate(name string) *XTemplate {
	if t.SubTemplates == nil {
		return nil
	}
	return t.SubTemplates[name]
}

// Execute will inject the Data into the template and creates the final string
func (t *XTemplate) Execute(data XDatasetDef) string {
	// Does data has a language ?
	if data != nil {
		lang, _ := data.Get("#")
		var language *XLanguage
		if lang != nil {
			language = lang.(*XLanguage)
		}
		stack := &XDatasetCollection{}
		stack.Push(data)
		return t.injector(stack, language)
	}
	return t.injector(nil, nil)
}

// injector will injects the data into this template
func (t *XTemplate) injector(datacol XDatasetCollectionDef, language *XLanguage) string {
	var injected []string
	if t.Root == nil {
		return "Error, no template.Root compiled"
	}
	for _, v := range *t.Root {
		switch v.ParamType {
		case MetaString: // included string from original code
			injected = append(injected, v.Data)
		case MetaComment:
			// nothing to do: comment ignored
		case MetaLanguage:
			if language != nil {
				injected = append(injected, language.Get(v.Data))
			}
		case MetaReference: // Reference
			// search for the subtemplate and reinject anything into it
			subt := t.GetTemplate(v.Data)
			if subt != nil {
				// if v.data is a substructure into data, then we stack the data and inject new stacked data

				substr := subt.injector(datacol, language)
				injected = append(injected, substr)
			}
		case MetaVariable:
			if datacol != nil {
				d, _ := datacol.GetDataString(v.Data)
				injected = append(injected, d)
			}
		case MetaRange: // Range (loop over subset) @@id:id@@
			// TODO(phil) subtemplate depends on first, last, loop # counter, key, condition etc
			xdata := strings.Split(v.Data, ":")
			subdataid := xdata[0]
			subtemplateid := xdata[0]
			if len(xdata) > 1 {
				subtemplateid = xdata[1]
			}

			subt := t.GetTemplate(subtemplateid)
			if subt != nil {
				// TODO(phil) We have to check the correct type of the collection

				if datacol != nil {
					cl, _ := datacol.GetCollection(subdataid)
					if cl != nil {
						for i := 0; i < cl.Count(); i++ {
							// if v.data is a substructure into data, then we stack the data and inject new stacked data
							dcl, _ := cl.Get(i)
							datacol.Push(dcl)
							substr := subt.injector(datacol, language)
							injected = append(injected, substr)
							// unstack extra data
							datacol.Pop()
						}
					}
				}
			}
		case MetaCondition:
			// TODO(phil) subtemplate depends on condition completion
			subt := t.GetTemplate(v.Data)
			value := searchConditionValue(v.Data, datacol)

			// if v.data is a substructure into data, then we stack the data and inject new stacked data
			if value != "" {
				substr := subt.injector(datacol, language)
				injected = append(injected, substr)
			}
		case MetaDump:
			if datacol != nil {
				dsubstr, _ := datacol.Get(0)
				if dsubstr != nil {
					substr := dsubstr.Stringify()
					injected = append(injected, substr)
				}
			}
		default:
			injected = append(injected, "THE METALANGUAGE FROM OUTERSPACE IS NOT SUPPORTED: "+fmt.Sprint(v.ParamType))
		}
	}
	// return the page string
	return strings.Join(injected, "")
}

// searchConditionValue will search one parameter value from the data
func searchConditionValue(id string, data XDatasetCollectionDef) string {
	// scan data for each dataset in order top to bottom
	if data == nil {
		return ""
	}
	v, _ := data.GetDataString(id)
	return v
}

/*
// buildValue will transform a data to a string
func buildValue(data interface{}) string {
	// if data is string, return data
	// if data is a type, return conversion
	// if data is a function, call the function then return the value
	return fmt.Sprint(data)
}
*/

// Print will creates the final string representing the code of the template
func (t *XTemplate) Print() string {
	return fmt.Sprint(t)
}
