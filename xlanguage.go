package xcore

import (
  "fmt"
  "os"
  "io/ioutil"
  "encoding/xml"
)

type XLanguageEntry struct {
  ID    string    `xml:"id,attr"`
  Entry string    `xml:",chardata"`
}

type XLanguage struct {
  Name string
  Entries []XLanguageEntry `xml:"entry"`
}

func NewLanguage() *XLanguage {
  return &XLanguage{}
}

func (l *XLanguage) LoadFile(file string) {
  xmlFile, _ := os.Open(file)
  defer xmlFile.Close()
  byteValue, _ := ioutil.ReadAll(xmlFile)
  xml.Unmarshal(byteValue, l)
}

func (l *XLanguage) Print() string {
  return fmt.Sprint(l)
}
