package xcore

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	//	"sync"
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
//
//	type of elements present in the template
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

// Clone will make a copy of the param
func (tp *XTemplateParam) Clone() *XTemplateParam {
	return &XTemplateParam{ParamType: tp.ParamType, Data: tp.Data}
}

// XTemplateData is an Array of all the parameters into the template
type XTemplateData []XTemplateParam

// XTemplate is the plain template structure
type XTemplate struct {
	Name         string
	Root         *XTemplateData
	SubTemplates map[string]*XTemplate
	// mutex        sync.RWMutex
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
			`(%)--(.*?)--%(\n|\r|\r\n|\n\r)?` + // index based 1

			// ==== LANGUAGE INJECTION
			`|(#)#([a-zA-Z0-9-_\.]+?)##` + // index based 4

			// ==== ELEMENTS
			`|(&)&([a-zA-Z0-9-_\=\>\:\|\.]+?)&&` + // index based 6
			`|(@)@([a-zA-Z0-9-_\=\>\:\|\.]+?)@@` + // index based 8
			`|(\?)\?([a-zA-Z0-9-_\=\>\:\|\.]+?)\?\?` + // index based 10
			`|(\!)\!([a-zA-Z0-9-_\=\>\:\|\.]+?)\!\!` + // index based 12
			`|(\{)\{([a-zA-Z0-9-_\=\>\:\|\.]+?)\}\}` + // index based 14

			// ==== NESTED ELEMENTS (SUB TEMPLATES)
			`|\[\[(\])\](\n|\r|\r\n|\n\r)?` + // index based 16
			`|(\[)\[([a-z0-9\|\.\-_]+?)\]\](\n|\r|\r\n|\n\r)?` // index based 18

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
		} else if matches[i][4] == "#" {
			param.ParamType = MetaLanguage // Language entry
			param.Data = matches[i][5]
		} else if matches[i][6] == "&" {
			param.ParamType = MetaReference // Reference to template
			param.Data = matches[i][7]
		} else if matches[i][8] == "@" {
			param.ParamType = MetaRange // Loop on data
			param.Data = matches[i][9]
		} else if matches[i][10] == "?" {
			param.ParamType = MetaCondition // Conditional on data
			param.Data = matches[i][11]
		} else if matches[i][12] == "!" {
			param.ParamType = MetaDump // Debug
			param.Data = matches[i][13]
		} else if matches[i][14] == "{" {
			param.ParamType = MetaVariable // Simple element
			param.Data = matches[i][15]
		} else if matches[i][18] == "[" {
			param.ParamType = MetaTemplateStart // Template start
			param.Data = matches[i][19]
		} else if matches[i][16] == "]" {
			param.ParamType = MetaTemplateEnd // Template end
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
		var language *XLanguage
		lang, _ := data.Get("#")
		if lang != nil {
			language, _ = lang.(*XLanguage) // language is nil if it-s not a *XLanguage
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
		case MetaReference: // Reference &&
			xid := strings.Split(v.Data, ":")
			if len(xid) == 3 {
				field := xid[1]
				prefix := xid[2]
				value, _ := datacol.GetDataString(field)
				subt := t.GetTemplate(prefix + value)
				if subt != nil {
					substr := subt.injector(datacol, language)
					injected = append(injected, substr)
				} else {
					subt := t.GetTemplate(prefix)
					if subt != nil {
						substr := subt.injector(datacol, language)
						injected = append(injected, substr)
					}
				}
			} else {
				template := ""
				if len(xid) >= 1 {
					template = xid[0]
				}
				subt := t.GetTemplate(template)
				if subt != nil {
					withds := false
					if len(xid) == 2 {
						dcl, _ := datacol.GetData(xid[1])
						ds, ok := dcl.(XDatasetDef)
						if ok {
							withds = true
							datacol.Push(ds)
						}
					}
					substr := subt.injector(datacol, language)
					if withds {
						datacol.Pop()
					}
					injected = append(injected, substr)
				}
			}
		case MetaVariable: // {{id>id>id...}}
			if datacol != nil {
				d, _ := datacol.GetDataString(v.Data)
				injected = append(injected, d)
			}
		case MetaRange: // Range (loop over subset) @@id:id@@
			xdata := strings.Split(v.Data, ":")
			subdataid := xdata[0]
			subtemplateid := xdata[0]
			if len(xdata) > 1 {
				subtemplateid = xdata[1]
			}

			subt := t.GetTemplate(subtemplateid)
			if subt != nil {
				if datacol != nil {
					cl, _ := datacol.GetCollection(subdataid)
					if cl != nil && cl.Count() > 0 {
						for i := 0; i < cl.Count(); i++ {
							var tmp *XTemplate
							tmp = t.GetTemplate(subtemplateid + ".key." + strconv.Itoa(i))
							//							if tmp == nil {
							//								tmp = t.GetTemplate(subtemplateid + ".field." + field + "." + value)
							//							}
							if tmp == nil && i == 0 {
								tmp = t.GetTemplate(subtemplateid + ".first")
							}
							if tmp == nil && i == cl.Count()-1 {
								tmp = t.GetTemplate(subtemplateid + ".last")
							}
							if tmp == nil && i%2 == 0 {
								tmp = t.GetTemplate(subtemplateid + ".even")
							}
							if tmp == nil {
								tmp = subt
							}
							dcl, _ := cl.Get(i)
							dcl.Set(".counter", i+1)
							datacol.Push(dcl)
							substr := tmp.injector(datacol, language)
							injected = append(injected, substr)
							// unstack extra data
							datacol.Pop()
						}
					} else {
						var tmp *XTemplate
						tmp = t.GetTemplate(subtemplateid + ".none")
						if tmp == nil {
							tmp = subt
						}
						substr := tmp.injector(datacol, language)
						injected = append(injected, substr)
					}
				}
			}
		case MetaCondition: //  ??id??
			xdata := strings.Split(v.Data, ":")
			subdataid := xdata[0]
			subtemplateid := xdata[0]
			if len(xdata) > 1 {
				subtemplateid = xdata[1]
			}

			subt := t.GetTemplate(subtemplateid)
			var value interface{}
			if datacol != nil {
				value, _ = datacol.GetData(subdataid)
			}
			if subt != nil && value != nil {
				withds := false
				svalue := ""
				ds, ok := value.(XDatasetDef)
				if ok {
					withds = true
					datacol.Push(ds)
					svalue = "array"
				} else {
					svalue = fmt.Sprint(value)
				}
				if svalue != "" {
					// subtemplate with .value?
					tmp := t.GetTemplate(subtemplateid + "." + svalue)
					if tmp != nil {
						subt = tmp
					}
					substr := subt.injector(datacol, language)
					injected = append(injected, substr)
				}
				if withds {
					datacol.Pop()
				}
			}
			if value == nil || fmt.Sprint(value) == "" {
				tmp := t.GetTemplate(subtemplateid + ".none")
				if tmp != nil {
					subt = tmp
				}
				if subt != nil {
					substr := subt.injector(datacol, language)
					injected = append(injected, substr)
				}
			}
		case MetaDump:
			if datacol != nil {
				if v.Data == "dump" || v.Data == "list" {
					dsubstr, _ := datacol.Get(0)
					if dsubstr != nil {
						substr := dsubstr.GoString()
						injected = append(injected, substr)
					}
				}
			}
		default:
			injected = append(injected, "THE METALANGUAGE FROM OUTERSPACE IS NOT SUPPORTED: "+fmt.Sprint(v.ParamType))
		}
	}
	// return the page string
	return strings.Join(injected, "")
}

// String will transform the XDataset into a readable string for humans
func (t *XTemplate) String() string {
	sdata := []string{}
	for _, val := range *t.Root {
		sdata = append(sdata, fmt.Sprintf("%v", val))
	}
	sort.Strings(sdata) // Lets be sure the print is always the same presentation
	return "xcore.XTemplate{" + strings.Join(sdata, " ") + "}"
}

// GoString will transform the XDataset into a readable string for humans with values indexes
func (t *XTemplate) GoString() string {
	sdata := []string{}
	for _, val := range *t.Root {
		sdata = append(sdata, fmt.Sprintf("%#v", val))
	}
	sort.Strings(sdata) // Lets be sure the print is always the same presentation
	return "#xcore.XTemplate{" + strings.Join(sdata, " ") + "}"
}

// Clone will make a full new copy of the template into a new memory space
func (t *XTemplate) Clone() *XTemplate {
	cloned := &XTemplate{Name: t.Name}
	if t.Root != nil {
		var newroot XTemplateData
		for _, td := range *t.Root {
			newroot = append(newroot, *td.Clone())
		}
		cloned.Root = &newroot
	}
	if t.SubTemplates != nil {
		newsubtemplates := map[string]*XTemplate{}
		for id, xt := range t.SubTemplates {
			newsubtemplates[id] = xt.Clone()
		}
		cloned.SubTemplates = newsubtemplates
	}

	return cloned
}
