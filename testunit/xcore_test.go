package main

import (
  "fmt"
  "testing"
  "github.com/webability-go/xcore"
)

/* TEST XLANGUAGE */
func TestXLanguage(t *testing.T) {
  // Test 1: creates from string
  lang := xcore.NewXLanguage()
  lang.LoadXML"
<?xml version="1.0" encoding="UTF-8"?>
<language id="language-demo" lang="en">
  <entry id="entry1">Welcome to</entry>
  <entry id="entry2">Xamboo</entry>
</language>
")
  
  


}








/* TEST XTEMPLATE */

func TestCommentParam(t *testing.T) {
  // Test 1: assign a simple parameter string with some comments
  tmpl := xcore.NewXTemplate()
  tmpl.LoadString("%-- starting comment --%Text %--with a [[]]comment--% here. Also an%----% empty comment %--ending comment--%")
  
//  fmt.Println(tmpl)
//  fmt.Println(tmpl.Root)
  
  result := tmpl.Execute(nil)
  
  fmt.Println("Result: ", result)

  // direct access
  if result != "Text  here. Also an empty comment " {
    t.Errorf("The comment has not been removed correctly")
  }
}

func TestLanguageParam(t *testing.T) {
  // Test 1: assign a simple parameter string with some comments
  tmpl := xcore.NewXTemplate()
  tmpl.LoadString("Test with ##some## ##languages## here")
  
  fmt.Println(tmpl)
  fmt.Println(tmpl.Root)
  
  data := make(map[string]interface{})
  data["language"] = NewXLanguage().Lo
  
  result := tmpl.Execute()
  
  fmt.Println("Result: ", result)

  // direct access
  if result != "Text  here. Also an empty comment " {
    t.Errorf("The comment has not been removed correctly")
  }
}


