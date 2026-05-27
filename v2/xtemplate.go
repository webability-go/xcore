package xcore

import (
	"errors"
	"fmt"
	"html"
	"io"
	"net/url"
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
  {{field|html}}    -- HTML-escape the value before injection
  {{field|js}}      -- JavaScript string-escape the value
  {{field|url}}     -- URL-query-encode the value
  {{field|raw}}     -- no escaping (default when no filter is given)
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

Thread safety note (fix 3.10):
  An XTemplate is safe for concurrent Execute / ExecuteSafe calls once it has been
  fully compiled (NewXTemplateFromString, NewXTemplateFromFile, or LoadString/LoadFile).
  Methods that mutate the template (AddTemplate, LoadFile, LoadString) must NOT be
  called concurrently with Execute or with each other. The recommended pattern is
  "compile once, execute many": build the template at start-up in a single goroutine,
  then share the read-only compiled value freely across goroutines.

Output escaping note (fix 3.1):
  By default {{field}} injects values without any escaping. When a template produces
  HTML and field values can come from untrusted sources, always use an output filter:
    {{name|html}}   -- safe for HTML text content / attributes
    {{cb|js}}       -- safe inside a JS string literal
    {{q|url}}       -- safe as a URL query parameter value
*/

// MetaString and other consts:
//
//	type of elements present in the template
const (
	MetaString  = 0 // a simple string to integrate into the code
	MetaComment = 1 // Comment, ignore it

	MetaLanguage  = 2 // one param of the URL parameters list, index-1 based
	MetaReference = 3 // reference to a sub-template
	MetaRange     = 4 // loop over a collection
	MetaCondition = 5 // conditional injection
	MetaDump      = 6 // debug dump
	MetaVariable  = 7 // simple variable injection

	MetaTemplateStart = 101 // Temporal nested box start tag
	MetaTemplateEnd   = 102 // Temporal nested box end tag

	MetaUnused = -1 // a "not used anymore" param to be freed
)

// maxTemplateDepth is the maximum allowed sub-template nesting depth.
// Recursive or mutually-referencing templates are aborted at this depth (fix 3.2 / 3.5).
const maxTemplateDepth = 64

// maxTemplateFileSize is the upper bound for template files loaded from disk.
// Files larger than this limit are rejected to prevent memory exhaustion (fix 3.6).
const maxTemplateFileSize = 10 * 1024 * 1024 // 10 MB

// XTemplateParam is a parameter definition into the template
type XTemplateParam struct {
	ParamType int
	Data      string
}

// Clone will make a copy of the param
func (tp *XTemplateParam) Clone() *XTemplateParam {
	return &XTemplateParam{ParamType: tp.ParamType, Data: tp.Data}
}

// XTemplateData is an Array of all the parameters into the template
type XTemplateData []XTemplateParam

// XTemplate is the plain template structure.
// See the package-level comment for thread-safety and output-escaping rules.
type XTemplate struct {
	Name         string
	Root         *XTemplateData
	SubTemplates map[string]*XTemplate
	Father       *XTemplate
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

// LoadFile will load a file into the template.
//
// WARNING (fix 3.7): the file path is passed directly to os.Open without sanitisation.
// If the path originates from user input, validate and restrict it to a known safe
// directory before calling this function to prevent path-traversal attacks.
func (t *XTemplate) LoadFile(file string) error {
	// 3.6: validate file size before reading to guard against DoS
	info, err := os.Stat(file)
	if err != nil {
		return err
	}
	if info.Size() > maxTemplateFileSize {
		return fmt.Errorf("xcore.XTemplate.LoadFile: file %q size %d exceeds limit %d", file, info.Size(), maxTemplateFileSize)
	}
	tFile, err := os.Open(file)
	if err != nil {
		return err
	}
	// 3.8: defer Close so the descriptor is always released even if ReadAll fails
	defer tFile.Close()
	// 3.8: io.ReadAll replaces deprecated ioutil.ReadAll
	data, err := io.ReadAll(tFile)
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
	code :=
		`(?s)` +
			`(%)--(.*?)--%(\n|\r|\r\n|\n\r)?` +
			`|(#)#([a-zA-Z0-9-_\.]+?)##` +
			`|(&)&([a-zA-Z0-9-_\=\>\:\|\.]+?)&&` +
			`|(@)@([a-zA-Z0-9-_\=\>\:\|\.]+?)@@` +
			`|(\?)\?([a-zA-Z0-9-_\=\>\:\|\.]+?)\?\?` +
			`|(\!)\!([a-zA-Z0-9-_\=\>\:\|\.]+?)\!\!` +
			`|(\{)\{([a-zA-Z0-9-_\=\>\:\|\.]+?)\}\}` +
			`|\[\[(\])\](\n|\r|\r\n|\n\r)?` +
			`|(\[)\[([a-z0-9\|\.\-_]+?)\]\](\n|\r|\r\n|\n\r)?`

	codex := regexp.MustCompile(code)
	indexes := codex.FindAllStringIndex(data, -1)
	matches := codex.FindAllStringSubmatch(data, -1)

	var compiled XTemplateData
	pointer := 0
	for i, x := range indexes {
		if pointer != x[0] {
			compiled = append(compiled, XTemplateParam{ParamType: MetaString, Data: data[pointer:x[0]]})
		}
		param := XTemplateParam{}
		if matches[i][1] == "%" {
			param.ParamType = MetaComment
			param.Data = matches[i][2]
		} else if matches[i][4] == "#" {
			param.ParamType = MetaLanguage
			param.Data = matches[i][5]
		} else if matches[i][6] == "&" {
			param.ParamType = MetaReference
			param.Data = matches[i][7]
		} else if matches[i][8] == "@" {
			param.ParamType = MetaRange
			param.Data = matches[i][9]
		} else if matches[i][10] == "?" {
			param.ParamType = MetaCondition
			param.Data = matches[i][11]
		} else if matches[i][12] == "!" {
			param.ParamType = MetaDump
			param.Data = matches[i][13]
		} else if matches[i][14] == "{" {
			param.ParamType = MetaVariable
			param.Data = matches[i][15]
		} else if matches[i][18] == "[" {
			param.ParamType = MetaTemplateStart
			param.Data = matches[i][19]
		} else if matches[i][16] == "]" {
			param.ParamType = MetaTemplateEnd
		} else {
			param.ParamType = MetaUnused
		}
		compiled = append(compiled, param)
		pointer = x[1]
	}
	if pointer != len(data) {
		compiled = append(compiled, XTemplateParam{ParamType: MetaString, Data: data[pointer:]})
	}

	// second pass: build sub-template hierarchy
	startpointers := []int{}
	subtemplates := []*XTemplate{}
	actualtemplate := t
	for i, x := range compiled {
		if x.ParamType == MetaTemplateStart {
			startpointers = append(startpointers, i)
			subtemplates = append(subtemplates, actualtemplate)
			actualtemplate = &XTemplate{Name: x.Data}
		} else if x.ParamType == MetaTemplateEnd {
			last := len(startpointers) - 1
			if last < 0 {
				return errors.New("Error: template mismatch Start/End")
			}
			startpointer := startpointers[last]
			startpointers = startpointers[:last]

			var subset XTemplateData
			for ptr := startpointer + 1; ptr < i; ptr++ {
				if compiled[ptr].ParamType != MetaUnused {
					subset = append(subset, compiled[ptr])
					compiled[ptr].ParamType = MetaUnused
				}
			}
			actualtemplate.Root = &subset

			uppertemplate := subtemplates[last]
			subtemplates = subtemplates[:last]

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
			actualtemplate = uppertemplate
			compiled[startpointer].ParamType = MetaUnused
			compiled[i].ParamType = MetaUnused
		}
	}
	if len(startpointers) > 0 {
		return errors.New("Error: template mismatch Start/End")
	}

	// last pass: compact array
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
	tmpl.Father = t
	t.SubTemplates[name] = tmpl
}

// GetTemplate gets a sub template by name, searching up the Father chain.
// Protected against cyclic Father pointers by a depth limit (fix 3.5).
func (t *XTemplate) GetTemplate(name string) *XTemplate {
	return t.getTemplateSafe(name, 0)
}

// getTemplateSafe is the cycle-safe implementation of GetTemplate (fix 3.5).
func (t *XTemplate) getTemplateSafe(name string, depth int) *XTemplate {
	if depth > maxTemplateDepth {
		return nil
	}
	if t.SubTemplates != nil {
		if tmp := t.SubTemplates[name]; tmp != nil {
			return tmp
		}
	}
	if t.Father != nil {
		return t.Father.getTemplateSafe(name, depth+1)
	}
	return nil
}

// Execute will inject the Data into the template and return the final string.
// If the template Root has not been compiled, Execute returns an empty string.
// Use ExecuteSafe to distinguish a genuine empty output from a compilation error (fix 3.9).
func (t *XTemplate) Execute(data XDatasetDef) string {
	result, _ := t.ExecuteSafe(data)
	return result
}

// ExecuteSafe will inject the Data into the template and return the final string
// together with any error. An error is returned only when the template Root has not
// been compiled (fix 3.9).
func (t *XTemplate) ExecuteSafe(data XDatasetDef) (string, error) {
	if t.Root == nil {
		return "", errors.New("xcore.XTemplate.ExecuteSafe: template Root has not been compiled")
	}
	if data != nil {
		var language *XLanguage
		lang, _ := data.Get("#")
		if lang != nil {
			language, _ = lang.(*XLanguage)
		}
		stack := &XDatasetCollection{}
		stack.Push(data)
		return t.injectorDepth(stack, language, 0), nil
	}
	return t.injectorDepth(nil, nil, 0), nil
}

// injectorDepth is the main injection engine with depth tracking (fix 3.2).
func (t *XTemplate) injectorDepth(datacol XDatasetCollectionDef, language *XLanguage, depth int) string {
	if t.Root == nil {
		return ""
	}
	if depth > maxTemplateDepth {
		return ""
	}
	var injected []string
	for _, v := range *t.Root {
		switch v.ParamType {
		case MetaString:
			injected = append(injected, v.Data)

		case MetaComment:
			// ignored

		case MetaLanguage:
			if language != nil {
				injected = append(injected, language.Get(v.Data))
			}

		case MetaReference:
			xid := strings.Split(v.Data, ":")
			if len(xid) == 3 {
				field := xid[1]
				prefix := xid[2]
				value, _ := datacol.GetDataString(field)
				subt := t.GetTemplate(prefix + value)
				if subt != nil {
					injected = append(injected, subt.injectorDepth(datacol, language, depth+1))
				} else {
					subt = t.GetTemplate(prefix)
					if subt != nil {
						injected = append(injected, subt.injectorDepth(datacol, language, depth+1))
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
					injected = append(injected, subt.injectorDepth(datacol, language, depth+1))
					if withds {
						datacol.Pop()
					}
				}
			}

		case MetaVariable: // {{path|filter}}
			if datacol != nil {
				// 3.1: optional output filter prevents XSS
				fieldPath, filter := splitFilter(v.Data)
				d, _ := datacol.GetDataString(fieldPath)
				injected = append(injected, applyOutputFilter(d, filter))
			}

		case MetaRange:
			xdata := strings.Split(v.Data, ":")
			subdataid := xdata[0]
			subtemplateid := xdata[0]
			if len(xdata) > 1 {
				subtemplateid = xdata[1]
			}
			subt := t.GetTemplate(subtemplateid)
			if subt != nil && datacol != nil {
				cl, _ := datacol.GetCollection(subdataid)
				if cl != nil && cl.Count() > 0 {
					for i := 0; i < cl.Count(); i++ {
						tmp := t.GetTemplate(subtemplateid + ".key." + strconv.Itoa(i))
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
						// 3.4: guard against nil entry from Get
						dcl, ok := cl.Get(i)
						if !ok || dcl == nil {
							continue
						}
						dcl.Set(".counter", i+1)
						datacol.Push(dcl)
						injected = append(injected, tmp.injectorDepth(datacol, language, depth+1))
						datacol.Pop()
					}
				} else {
					tmp := t.GetTemplate(subtemplateid + ".none")
					if tmp == nil {
						tmp = subt
					}
					injected = append(injected, tmp.injectorDepth(datacol, language, depth+1))
				}
			}

		case MetaCondition:
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
					tmp := t.GetTemplate(subtemplateid + "." + svalue)
					if tmp != nil {
						subt = tmp
					}
					injected = append(injected, subt.injectorDepth(datacol, language, depth+1))
				}
				if withds {
					datacol.Pop()
				}
			}
			// value absent or empty: use .none; bug fix: use tmp not subt
			if value == nil || fmt.Sprint(value) == "" {
				tmp := t.GetTemplate(subtemplateid + ".none")
				if tmp == nil {
					tmp = subt
				}
				if tmp != nil {
					injected = append(injected, tmp.injectorDepth(datacol, language, depth+1))
				}
			}

		case MetaDump:
			if datacol != nil && (v.Data == "dump" || v.Data == "list") {
				dsubstr, _ := datacol.Get(0)
				if dsubstr != nil {
					injected = append(injected, dsubstr.GoString())
				}
			}

		default:
			injected = append(injected, "THE METALANGUAGE FROM OUTERSPACE IS NOT SUPPORTED: "+fmt.Sprint(v.ParamType))
		}
	}
	return strings.Join(injected, "")
}

// String will transform the XTemplate into a readable string for humans
func (t *XTemplate) String() string {
	if t.Root == nil {
		return "xcore.XTemplate{}"
	}
	sdata := []string{}
	for _, val := range *t.Root {
		sdata = append(sdata, fmt.Sprintf("%v", val))
	}
	sort.Strings(sdata)
	return "xcore.XTemplate{" + strings.Join(sdata, " ") + "}"
}

// GoString will transform the XTemplate into a readable string for humans with value indexes
func (t *XTemplate) GoString() string {
	if t.Root == nil {
		return "#xcore.XTemplate{}"
	}
	sdata := []string{}
	for _, val := range *t.Root {
		sdata = append(sdata, fmt.Sprintf("%#v", val))
	}
	sort.Strings(sdata)
	return "#xcore.XTemplate{" + strings.Join(sdata, " ") + "}"
}

// Clone will make a full new copy of the template into a new memory space.
// Mutations to the clone do not affect the original and vice versa.
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
			// 3.3: set Father on the CLONED sub, not on the original xt
			clonedSub := xt.Clone()
			clonedSub.Father = cloned
			newsubtemplates[id] = clonedSub
		}
		cloned.SubTemplates = newsubtemplates
	}
	cloned.Father = t.Father
	return cloned
}

// ---------------------------------------------------------------------------
// Output filter helpers (fix 3.1)
// ---------------------------------------------------------------------------

// splitFilter separates a field path from an optional output-filter suffix.
//
//	"name"      -> ("name", "")
//	"name|html" -> ("name", "html")
//	"a>b|url"   -> ("a>b", "url")
func splitFilter(data string) (fieldPath, filter string) {
	if idx := strings.LastIndex(data, "|"); idx >= 0 {
		return data[:idx], data[idx+1:]
	}
	return data, ""
}

// applyOutputFilter applies the named output filter to value.
// Supported filters: "html", "js", "url" / "urlquery", "raw" / "" (no-op).
func applyOutputFilter(value, filter string) string {
	switch filter {
	case "html":
		return html.EscapeString(value)
	case "js":
		return escapeJS(value)
	case "url", "urlquery":
		return url.QueryEscape(value)
	default:
		return value
	}
}

// escapeJS escapes s for safe embedding inside a JavaScript string literal.
// It neutralises HTML meta-characters (<, >, &) to prevent injection when
// the JS literal is inside an HTML document, and escapes U+2028 / U+2029
// which some JS engines treat as line terminators.
func escapeJS(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 8)
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		case '\'':
			b.WriteString(`\'`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		case '<':
			b.WriteString(`\x3C`)
		case '>':
			b.WriteString(`\x3E`)
		case '&':
			b.WriteString(`\x26`)
		case ' ': // Unicode line separator
			b.WriteString(` `)
		case ' ': // Unicode paragraph separator
			b.WriteString(` `)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
