package xcore

import (
  "fmt"
  "os"
  "io/ioutil"
  "encoding/xml"
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

Text format is:
ENTRYNAME=ENTRYVALUE
ENTRYNAME=ENTRYVALUE
where ENTRYNAME is the ID of the entry, for example "greating", "yourname", "submitbutton"
      ENTRYVALUE is the text for your entry, for example "Hello", "You are:", "Save" if your table is in english

There is no name of table or language in this format (you "know" what you are loading)

The advantage to use XML forma is to have more control over your language, and eventyally add attributes into your entry, 
for instance translated="yes/no", verified="yes/no", and any other data that your system could insert

There are 4 functions to create the language from a file or string, and 2 functions to get/set values

- NewXLanguage: to create an empty XLanguage structure
- NewXLanguageFromXMLFile: 
- NewXLanguageFromXMLString: 
- NewXLanguageFromFile: 
- NewXLanguageFromString: 

*/

// OFicial XLanguage structure for the user
type XLanguage struct {
  Name string
  Language string
  Entries map[string]string
}

/* NewXLanguage:
   Creates an empty Language structure with a name and a language
*/
func NewXLanguage(name string, lang string) *XLanguage {
  return &XLanguage{Name: name, Language: lang, Entries: make(map[string]string) }
}

/* NewXLanguageFromXMLFile:
   Creates an XLanguage structure with the data into the XML file
   Returns nil if there is an error
*/
func NewXLanguageFromXMLFile(file string) (*XLanguage, error) {
  lang := &XLanguage{Entries: make(map[string]string)}
  err := lang.LoadXMLFile(file)
  if err != nil {
    return nil, err
  }
  return lang, nil
}

func NewXLanguageFromXMLString(xml string) (*XLanguage, error) {
  lang := &XLanguage{Entries: make(map[string]string)}
  err := lang.LoadXMLString(xml)
  if err != nil {
    return nil, err
  }
  return lang, nil
}

func NewXLanguageFromFile(file string) (*XLanguage, error) {
  l := &XLanguage{Entries: make(map[string]string)}
  err := l.LoadFile(file)
  if err != nil {
    return nil, err
  }
  return l, nil
}

func NewXLanguageFromString(data string) (*XLanguage, error) {
  l := &XLanguage{Entries: make(map[string]string)}
  err := l.LoadString(data)
  if err != nil {
    return nil, err
  }
  return l, nil
}

/* LoadXMLFile:
   Loads a language from an XML file and replace the content of the XLanguage structure with the new data
*/
func (l *XLanguage)LoadXMLFile(file string) error {
  xmlFile, err := os.Open(file)
  if err != nil { return err }
  data, err := ioutil.ReadAll(xmlFile)
  if err != nil { return err }
  err = xmlFile.Close()
  if err != nil { return err }
  return l.LoadXMLString(string(data))
}

func (l *XLanguage)LoadXMLString(data string) error {
  // Temporal structures for XML loading
  type xlanguageentrytemp struct {
    ID    string    `xml:"id,attr"`
    Entry string    `xml:",chardata"`
  }

  type xlanguagetemp struct {
    Name     string  `xml:"id,attr"`
    Language string  `xml:"lang,attr"`
    Entries  []xlanguageentrytemp `xml:"entry"`
  }

  // Unmarshal
  temp := &xlanguagetemp{}
  err := xml.Unmarshal([]byte(data), temp)
  if err != nil { return err }
  
  // Scan to our XLanguage Object
  l.Name = temp.Name
  l.Language = temp.Language
  for _, e := range temp.Entries {
    l.Entries[e.ID] = e.Entry
  }
  return nil
}

func (l *XLanguage)LoadFile(file string) error {
  xmlFile, _ := os.Open(file)
  defer xmlFile.Close()
  byteValue, _ := ioutil.ReadAll(xmlFile)
  
  fmt.Println(byteValue)
  
  return nil
}

func (l *XLanguage)LoadString(data string) error {
  return nil
}

func (l *XLanguage)SetName(name string) {
  l.Name = name
}

func (l *XLanguage)SetLanguage(lang string) {
  l.Language = lang
}

func (l *XLanguage)GetName() string {
  return l.Name
}

func (l *XLanguage)GetLanguage() string {
  return l.Language
}

func (l *XLanguage)Set(entry string, value string) {
  l.Entries[entry] = value
}

func (l *XLanguage)Get(entry string) string {
  v, ok := l.Entries[entry]
  if ok { return v }
  return ""
}

func (l *XLanguage)Del(entry string) {
  delete(l.Entries, entry)
}

func (l *XLanguage) Print() string {
  return fmt.Sprint(l)
}
