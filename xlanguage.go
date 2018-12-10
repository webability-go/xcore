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
  return &XLanguage{Name: name, Language: lang}
}

/* NewXLanguageFromXMLFile:
   Creates an XLanguage structure with the data into the XML file
   Returns nil if there is an error
*/
func NewXLanguageFromXMLFile(file string) (*XLanguage, error) {
  lang := &XLanguage{}
  err := lang.LoadXMLFile(file)
  if err != nil {
    return nil, err
  }
  return lang, nil
}

func NewXLanguageFromXML(xml string) (*XLanguage, error) {
  lang := &XLanguage{}
  err := lang.LoadXML(xml)
  if err != nil {
    return nil, err
  }
  return lang, nil
}

func NewXLanguageFromFile(data string) *XLanguage {
  return &XLanguage{}.LoadString(data)
}


// Temporal structures for XML loading
type xlanguageentrytemp struct {
  ID    string    `xml:"id,attr"`
  Entry string    `xml:",chardata"`
}

type xlanguagetemp struct {
  Entries []XLanguageEntry `xml:"entry"`
}

/* LoadXMLFile:
   Loads a language from an XML file and replace the content of the XLanguage structure with the new data
*/
func (l *XLanguage)LoadXMLFile(file string) error {
  xmlFile, _ := os.Open(file)
  defer xmlFile.Close()
  byteValue, _ := ioutil.ReadAll(xmlFile)
  xml.Unmarshal(byteValue, l)
  // Build the map entry
  
}

func (l *XLanguage)LoadString(data string) {
  xml.Unmarshal(data, l)
}

func (l *XLanguage) LoadXML(data string) {
  xml.Unmarshal(data, l)
}

func (l *XLanguage)Get(entry string) string {
  for _, v := range l.Entries {
    if v.ID == entry {
      return v.Entry
  }
  return ""
}

func (l *XLanguage)Set(entry string, value string) {
  e := &XLanguageEntry{ID: entry, Entry: value}
  l.Entries := append(l.Entries, e)
}

func (l *XLanguage) Print() string {
  return fmt.Sprint(l)
}


