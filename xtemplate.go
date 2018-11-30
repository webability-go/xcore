package xcore

import (
  "fmt"
  "os"
  "io/ioutil"
)

type XTemplate struct {
  Name string
  Content string
}

func NewTemplate() *XTemplate {
  return &XTemplate{}
}

func (t *XTemplate) LoadFile(file string) {
  xmlFile, _ := os.Open(file)
  defer xmlFile.Close()
  byteValue, _ := ioutil.ReadAll(xmlFile)
  t.Content = string(byteValue)
}

func (t *XTemplate) Print() string {
  return fmt.Sprint(t)
}
