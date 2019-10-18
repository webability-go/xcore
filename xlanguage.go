package xcore

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

/*
The XLanguage table of text entries can be loaded from XML file, XML string or normal file or string.
XML Format is:
<?xml version="1.0" encoding="UTF-8"?>
<language id="NAMEOFLANGUAGE" lang="LG">
  <entry id="ENTRYNAME">ENTRYVALUE</entry>
  <entry id="ENTRYNAME">ENTRYVALUE</entry>
</language>
where NAMEOFLANGUAGE is the name of your table entry, for example "homepate", "user_report", etc
      LG is the ISO-3369 2 letters language ID, for example "es" for spanish, "en" for english
      ENTRYNAME is the ID of the entry, for example "greating", "yourname", "submitbutton"
      ENTRYVALUE is the text for your entry, for example "Hello", "You are:", "Save" if your table is in english

Flat Text format is:
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

*/

// XLanguage is the oficial structure for the user
type XLanguage struct {
	Name     string
	Language string
	Entries  map[string]string
}

// NewXLanguage will create an empty Language structure with a name and a language
func NewXLanguage(name string, lang string) *XLanguage {
	return &XLanguage{Name: name, Language: lang, Entries: make(map[string]string)}
}

// NewXLanguageFromXMLFile will create an XLanguage structure with the data into the XML file
//   Returns nil if there is an error
func NewXLanguageFromXMLFile(file string) (*XLanguage, error) {
	lang := &XLanguage{Entries: make(map[string]string)}
	err := lang.LoadXMLFile(file)
	if err != nil {
		return nil, err
	}
	return lang, nil
}

// NewXLanguageFromXMLString will create an XLanguage structure with the data into the XML String
//   Returns nil if there is an error
func NewXLanguageFromXMLString(xml string) (*XLanguage, error) {
	lang := &XLanguage{Entries: make(map[string]string)}
	err := lang.LoadXMLString(xml)
	if err != nil {
		return nil, err
	}
	return lang, nil
}

// NewXLanguageFromFile will create an XLanguage structure with the data into the text file
//   Returns nil if there is an error
func NewXLanguageFromFile(file string) (*XLanguage, error) {
	l := &XLanguage{Entries: make(map[string]string)}
	err := l.LoadFile(file)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// NewXLanguageFromString will create an XLanguage structure with the data into the string
//   Returns nil if there is an error
func NewXLanguageFromString(data string) (*XLanguage, error) {
	l := &XLanguage{Entries: make(map[string]string)}
	err := l.LoadString(data)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// LoadXMLFile will Load a language from an XML file and replace the content of the XLanguage structure with the new data
//   Returns nil if there is an error
func (l *XLanguage) LoadXMLFile(file string) error {
	xmlFile, err := os.Open(file)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return err
	}
	err = xmlFile.Close()
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
		Id    string `xml:"id,attr"`
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
	l.Language = temp.Language
	for _, e := range temp.Entries {
		l.Entries[e.Id] = e.Entry
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
		l.Entries[key] = value
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
func (l *XLanguage) SetLanguage(lang string) {
	l.Language = lang
}

// GetName will return the name of the language table
func (l *XLanguage) GetName() string {
	return l.Name
}

// GetLanguage will return the language of the language table
func (l *XLanguage) GetLanguage() string {
	return l.Language
}

// Set will add an entry id-value into the language table
func (l *XLanguage) Set(entry string, value string) {
	l.Entries[entry] = value
}

// Get will read an entry id-value from the language table
func (l *XLanguage) Get(entry string) string {
	v, ok := l.Entries[entry]
	if ok {
		return v
	}
	return ""
}

// Del will remove an entry id-value from the language table
func (l *XLanguage) Del(entry string) {
	delete(l.Entries, entry)
}

// Print will get the language table as a readable string
func (l *XLanguage) Print() string {
	return fmt.Sprint(l)
}
