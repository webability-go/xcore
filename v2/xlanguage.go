package xcore

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"sync"

	"golang.org/x/text/language"
)

// XLanguage is the oficial structure for the user
type XLanguage struct {
	Name     string
	Language language.Tag
	entries  map[string]string
	mutex    sync.RWMutex
}

// NewXLanguage will create an empty Language structure with a name and a language
func NewXLanguage(name string, lang language.Tag) *XLanguage {
	return &XLanguage{Name: name, Language: lang, entries: make(map[string]string)}
}

// NewXLanguageFromXMLFile will create an XLanguage structure with the data into the XML file
//   Returns nil if there is an error
func NewXLanguageFromXMLFile(file string) (*XLanguage, error) {
	lang := &XLanguage{entries: make(map[string]string)}
	err := lang.LoadXMLFile(file)
	if err != nil {
		return nil, err
	}
	return lang, nil
}

// NewXLanguageFromXMLString will create an XLanguage structure with the data into the XML String
//   Returns nil if there is an error
func NewXLanguageFromXMLString(xml string) (*XLanguage, error) {
	lang := &XLanguage{entries: make(map[string]string)}
	err := lang.LoadXMLString(xml)
	if err != nil {
		return nil, err
	}
	return lang, nil
}

// NewXLanguageFromFile will create an XLanguage structure with the data into the text file
//   Returns nil if there is an error
func NewXLanguageFromFile(file string) (*XLanguage, error) {
	l := &XLanguage{entries: make(map[string]string)}
	err := l.LoadFile(file)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// NewXLanguageFromString will create an XLanguage structure with the data into the string
//   Returns nil if there is an error
func NewXLanguageFromString(data string) (*XLanguage, error) {
	l := &XLanguage{entries: make(map[string]string)}
	err := l.LoadString(data)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// LoadXMLFile will Load a language from an XML file and replace the content of the XLanguage structure with the new data
//   Returns nil if there is an error
func (l *XLanguage) LoadXMLFile(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return l.LoadXMLString(string(data))
}

// LoadXMLString will Load a language from an XML file and replace the content of the XLanguage structure with the new data
//   Returns nil if there is an error
func (l *XLanguage) LoadXMLString(data string) error {
	// Temporal structures for XML loading
	type xentry struct {
		ID    string `xml:"id,attr"`
		Entry string `xml:",chardata"`
	}

	type xlang struct {
		Name     string   `xml:"id,attr"`
		Language string   `xml:"lang,attr"`
		Entries  []xentry `xml:"entry"`
	}

	// Unmarshal
	temp := &xlang{}
	err := xml.Unmarshal([]byte(data), temp)

	if err != nil {
		return err
	}

	// Scan to our XLanguage Object
	l.Name = temp.Name
	l.Language, _ = language.Parse(temp.Language)
	l.mutex.Lock()
	defer l.mutex.Unlock()
	for _, e := range temp.Entries {
		l.entries[e.ID] = e.Entry
	}
	return nil
}

// LoadFile will Load a language from a file and replace the content of the XLanguage structure with the new data
//   Returns nil if there is an error
func (l *XLanguage) LoadFile(file string) error {
	flatFile, err := os.Open(file)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(flatFile)
	if err != nil {
		return err
	}
	err = flatFile.Close()
	if err != nil {
		return err
	}
	return l.LoadString(string(data))
}

// LoadString will Load a language from a string and replace the content of the XLanguage structure with the new data
//   Returns nil if there is an error
func (l *XLanguage) LoadString(data string) error {
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

		// we capture the value if it exists. If not, the key entry is initialized with a nil value
		value := ""
		if len(line) > posequal {
			value = strings.TrimSpace(line[posequal+1:])
		}
		l.mutex.Lock()
		defer l.mutex.Unlock()
		l.entries[key] = value
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// SetName will set the name of the language table
func (l *XLanguage) SetName(name string) {
	l.Name = name
}

// SetLanguage will set the language ISO code (2 letters) of the language table
func (l *XLanguage) SetLanguage(lang language.Tag) {
	l.Language = lang
}

// GetName will return the name of the language table
func (l *XLanguage) GetName() string {
	return l.Name
}

// GetLanguage will return the language of the language table
func (l *XLanguage) GetLanguage() language.Tag {
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

// Del will remove an entry id-value from the language table
func (l *XLanguage) Del(entry string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	delete(l.entries, entry)
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

// String will transform the XDataset into a readable string for humans
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

// GoString will transform the XDataset into a readable string for humans
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
