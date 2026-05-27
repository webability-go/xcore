package xcore

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"sync"

	"golang.org/x/text/language"
)

// maxLanguageFileSize is the upper bound accepted for language files loaded from disk or raw strings.
// Files or strings larger than this limit are rejected to prevent memory exhaustion (DoS).
const maxLanguageFileSize = 10 * 1024 * 1024 // 10 MB

// XLanguage is the oficial structure for the user
type XLanguage struct {
	Name     string
	Language language.Tag
	entries  map[string]string
	status   map[string]string
	mutex    sync.RWMutex
}

// NewXLanguage will create an empty Language structure with a name and a language
func NewXLanguage(name string, lang language.Tag) *XLanguage {
	return &XLanguage{Name: name, Language: lang, entries: make(map[string]string), status: make(map[string]string)}
}

// NewXLanguageFromXMLFile will create an XLanguage structure with the data into the XML file
//
//	Returns nil if there is an error
func NewXLanguageFromXMLFile(file string) (*XLanguage, error) {
	lang := &XLanguage{entries: make(map[string]string), status: make(map[string]string)}
	err := lang.LoadXMLFile(file)
	if err != nil {
		return nil, err
	}
	return lang, nil
}

// NewXLanguageFromXMLString will create an XLanguage structure with the data into the XML String
//
//	Returns nil if there is an error
func NewXLanguageFromXMLString(xmldata string) (*XLanguage, error) {
	lang := &XLanguage{entries: make(map[string]string), status: make(map[string]string)}
	err := lang.LoadXMLString(xmldata)
	if err != nil {
		return nil, err
	}
	return lang, nil
}

// NewXLanguageFromFile will create an XLanguage structure with the data into the text file
//
//	Returns nil if there is an error
func NewXLanguageFromFile(file string) (*XLanguage, error) {
	l := &XLanguage{entries: make(map[string]string), status: make(map[string]string)}
	err := l.LoadFile(file)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// NewXLanguageFromString will create an XLanguage structure with the data into the string
//
//	Returns nil if there is an error
func NewXLanguageFromString(data string) (*XLanguage, error) {
	l := &XLanguage{entries: make(map[string]string), status: make(map[string]string)}
	err := l.LoadString(data)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// LoadXMLFile will Load a language from an XML file and replace the content of the XLanguage structure with the new data
//
//	Returns nil if there is an error
func (l *XLanguage) LoadXMLFile(file string) error {
	// 2.5: validate file size before reading to guard against DoS via oversized files
	info, err := os.Stat(file)
	if err != nil {
		return err
	}
	if info.Size() > maxLanguageFileSize {
		return fmt.Errorf("xcore.XLanguage.LoadXMLFile: file %q size %d exceeds limit %d", file, info.Size(), maxLanguageFileSize)
	}
	// 2.7: os.ReadFile replaces deprecated ioutil.ReadFile
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return l.LoadXMLString(string(data))
}

// LoadXMLString will Load a language from an XML string and replace the content of the XLanguage structure with the new data
//
//	Returns nil if there is an error
func (l *XLanguage) LoadXMLString(data string) error {
	// 2.5: guard against oversized strings
	if len(data) > maxLanguageFileSize {
		return fmt.Errorf("xcore.XLanguage.LoadXMLString: data size %d exceeds limit %d", len(data), maxLanguageFileSize)
	}

	// Temporal structures for XML loading
	type xentry struct {
		ID     string `xml:"id,attr"`
		Entry  string `xml:",chardata"`
		Status string `xml:"status,attr"`
	}

	type xlang struct {
		Name     string   `xml:"id,attr"`
		Language string   `xml:"lang,attr"`
		Entries  []xentry `xml:"entry"`
	}

	// 2.6: use an explicit Decoder with Strict = true (rejects malformed XML)
	temp := &xlang{}
	decoder := xml.NewDecoder(strings.NewReader(data))
	decoder.Strict = true
	if err := decoder.Decode(temp); err != nil {
		return err
	}

	// 2.2: surface language.Parse errors instead of silencing them
	lang, err := language.Parse(temp.Language)
	if err != nil {
		log.Printf("xcore.XLanguage.LoadXMLString: invalid language tag %q: %v", temp.Language, err)
		lang = language.Und
	}

	// 2.1: protect Name, Language and all entry writes under a single Lock
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.Name = temp.Name
	l.Language = lang
	for _, e := range temp.Entries {
		l.entries[e.ID] = e.Entry
		l.status[e.ID] = e.Status
	}
	return nil
}

// LoadFile will Load a language from a file and replace the content of the XLanguage structure with the new data
//
//	Returns nil if there is an error
func (l *XLanguage) LoadFile(file string) error {
	// 2.5: validate file size before reading
	info, err := os.Stat(file)
	if err != nil {
		return err
	}
	if info.Size() > maxLanguageFileSize {
		return fmt.Errorf("xcore.XLanguage.LoadFile: file %q size %d exceeds limit %d", file, info.Size(), maxLanguageFileSize)
	}
	flatFile, err := os.Open(file)
	if err != nil {
		return err
	}
	// 2.8: defer Close so the descriptor is always released, even on ReadAll error
	defer flatFile.Close()
	// 2.7: io.ReadAll replaces deprecated ioutil.ReadAll
	data, err := io.ReadAll(flatFile)
	if err != nil {
		return err
	}
	return l.LoadString(string(data))
}

// LoadString will Load a language from a string and replace the content of the XLanguage structure with the new data
//
//	Returns nil if there is an error
func (l *XLanguage) LoadString(data string) error {
	// 2.4: parse all lines into a local map first; apply a single Lock at the end so that
	// concurrent readers never observe a half-loaded state.
	local := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		posequal := strings.Index(line, "=")

		// we ignore empty and comments lines, no key=value lines too
		if len(line) == 0 || line[0] == '#' || line[0] == ';' || posequal < 0 {
			continue
		}

		// we separate the key. if there is no key, we ignore the data
		key := strings.TrimSpace(line[:posequal])
		if len(key) == 0 {
			continue
		}

		// we capture the value if it exists. If not, the key entry is initialized with an empty value
		value := ""
		if len(line) > posequal {
			value = strings.TrimSpace(line[posequal+1:])
		}
		local[key] = value
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	// Single lock for the entire update
	l.mutex.Lock()
	defer l.mutex.Unlock()
	for k, v := range local {
		l.entries[k] = v
	}
	return nil
}

// SetName will set the name of the language table
func (l *XLanguage) SetName(name string) {
	// 2.1: protect Name with mutex
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.Name = name
}

// SetLanguage will set the language ISO code (2 letters) of the language table
func (l *XLanguage) SetLanguage(lang language.Tag) {
	// 2.1: protect Language with mutex
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.Language = lang
}

// GetName will return the name of the language table
func (l *XLanguage) GetName() string {
	// 2.1: protect Name with RLock
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.Name
}

// GetLanguage will return the language of the language table
func (l *XLanguage) GetLanguage() language.Tag {
	// 2.1: protect Language with RLock
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.Language
}

// Set will add an entry id-value into the language table
func (l *XLanguage) Set(entry string, value string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.entries[entry] = value
}

// Get will read an entry id-value from the language table
func (l *XLanguage) Get(entry string) string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	v, ok := l.entries[entry]
	if ok {
		return v
	}
	return ""
}

// SetStatus will add a status value for an entry into the language table
func (l *XLanguage) SetStatus(entry string, value string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.status[entry] = value
}

// GetStatus will read a status value for an entry from the language table
func (l *XLanguage) GetStatus(entry string) string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	v, ok := l.status[entry]
	if ok {
		return v
	}
	return ""
}

// Del will remove an entry id-value from the language table
func (l *XLanguage) Del(entry string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	delete(l.entries, entry)
	delete(l.status, entry)
}

// GetEntries will return a COPY of the key-values pairs of the language.
func (l *XLanguage) GetEntries() map[string]string {
	clone := map[string]string{}
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	for k, v := range l.entries {
		clone[k] = v
	}
	return clone
}

// GetXML will generate the XML to save the file.
// All attribute values and text content are properly escaped via encoding/xml,
// preventing XML injection even when keys, statuses, or values contain
// special characters such as ", <, >, & or the CDATA terminator ]]>.
func (l *XLanguage) GetXML() string {
	// 2.3: use struct-based marshaling so encoding/xml handles all escaping
	type xentry struct {
		ID     string `xml:"id,attr"`
		Status string `xml:"status,attr"`
		Entry  string `xml:",chardata"`
	}
	type xlang struct {
		XMLName  xml.Name `xml:"language"`
		ID       string   `xml:"id,attr"`
		Language string   `xml:"lang,attr"`
		Entries  []xentry `xml:"entry"`
	}

	// Collect under RLock; release before marshaling (no lock held during allocation)
	l.mutex.RLock()
	keys := make([]string, 0, len(l.entries))
	for k := range l.entries {
		keys = append(keys, k)
	}
	sort.Strings(keys) // deterministic output
	entries := make([]xentry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, xentry{ID: k, Status: l.status[k], Entry: l.entries[k]})
	}
	name := l.Name
	langStr := fmt.Sprint(l.Language)
	l.mutex.RUnlock()

	doc := xlang{
		ID:       name,
		Language: langStr,
		Entries:  entries,
	}
	data, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return ""
	}
	return xml.Header + string(data)
}

// String will transform the XLanguage into a readable string for humans
func (l *XLanguage) String() string {
	sdata := []string{}
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	for key, val := range l.entries {
		sdata = append(sdata, key+":"+fmt.Sprintf("%v", val))
	}
	sort.Strings(sdata) // Lets be sure the print is always the same presentation
	return "xcore.XLanguage{" + strings.Join(sdata, " ") + "}"
}

// GoString will transform the XLanguage into a readable string for humans
func (l *XLanguage) GoString() string {
	sdata := []string{}
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	for key, val := range l.entries {
		sdata = append(sdata, key+":"+fmt.Sprintf("%#v", val))
	}
	sort.Strings(sdata) // Lets be sure the print is always the same presentation
	return "#xcore.XLanguage{" + strings.Join(sdata, " ") + "}"
}
