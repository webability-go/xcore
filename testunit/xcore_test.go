package main

import (
  "fmt"
  "testing"
  "github.com/webability-go/xcore"
)

/* TEST XLANGUAGE */
func TestXLanguage(t *testing.T) {
  // Test 1: creates from empty
  lang := xcore.NewXLanguage("mytranslation", "es")
  fmt.Println(lang)

  // Test 2: set/get
  lang.Set("greatings", "Hola, mundo")
  lang.Set("submit", "Guardar")
  
  tr := lang.Get("greatings")

  fmt.Println("Greatings: ", tr)

  lang.Del("greatings")
  tr2 := lang.Get("greatings")

  fmt.Println("Greatings after delete: ", tr2)

  if tr != "Hola, mundo" || tr2 != "" {
    t.Errorf("Set/Get/Del is not working correctly")
  }
  
  // Test 3: creates from XML string
  lang, err := xcore.NewXLanguageFromXMLString(`
<?xml version="1.0" encoding="UTF-8"?>
<language id="language-demo" lang="en">
  <entry id="entry1">Welcome to</entry>
  <entry id="entry2">XCore</entry>
</language>
`)
  
  fmt.Println(lang, err)
  
  tr = lang.Get("entry1")
  if tr != "Welcome to" {
    t.Errorf("NewXLanguageFromXMLString is not working correctly")
  }

  // Test 4: creates from Flat string
  lang, err = xcore.NewXLanguageFromString(`
entry1=Welcome to
entry2=XCore
`)
  
  fmt.Println(lang, err)
  
  tr = lang.Get("entry2")
  if tr != "XCore" {
    t.Errorf("NewXLanguageFromString is not working correctly")
  }
}







/* TEST XTEMPLATE */
/*
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
//  data["language"] = NewXLanguage().Lo
  
  result := tmpl.Execute(data)
  
  fmt.Println("Result: ", result)

  // direct access
  if result != "Text  here. Also an empty comment " {
    t.Errorf("The comment has not been removed correctly")
  }
}

*/
