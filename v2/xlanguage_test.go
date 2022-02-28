package xcore

import (
	"fmt"
	"io/ioutil"
	"testing"

	"golang.org/x/text/language"
)

func ExampleNewXLanguage() {

	langES := NewXLanguage("messages", language.Spanish)
	langEN := NewXLanguage("messages", language.English)
	langFR := NewXLanguage("messages", language.French)

	// You can load this from your system files for instance and keep your translation tables apart
	langES.Set("panicerror", "Error crítico del sistema")
	langEN.Set("panicerror", "System panic error")
	langFR.Set("panicerror", "Erreur grave dans le système")

	// pick a random general system language (for instance the client's OS language)
	lang := langFR

	// launch a panic errors
	fmt.Println("Launch a panic error message in the selected language:", lang.Get("panicerror"))
	// Output:
	// Launch a panic error message in the selected language: Erreur grave dans le système
}

func ExampleNewXLanguageFromXMLFile() {

	langES, _ := NewXLanguageFromXMLFile("./testunit/a.es.language")
	langEN, _ := NewXLanguageFromXMLFile("./testunit/a.en.language")
	langFR, _ := NewXLanguageFromXMLFile("./testunit/a.fr.language")

	// pick a random general system language (for instance the client's OS language)
	lang := langES

	// Say hello main language
	fmt.Println(lang.Get("entry1") + " " + lang.Get("entry2"))
	// Say hello in other languages
	fmt.Println(langEN.Get("entry1") + " " + langEN.Get("entry2"))
	fmt.Println(langFR.Get("entry1") + " " + langFR.Get("entry2"))
	// Output:

	// Bienvenido a XCore
	// Welcome to XCore
	// Bienvenue à XCore
}

func TestXLanguage(t *testing.T) {

	// Those strings are the same in errors.es.txt and errors.es.xml to make the tests
	teststrings := map[string]string{
		"panicerror":  "Error crítico del sistema",
		"systemerror": "Error del sistema",
		"fileerror":   "Error leyendo el archivo",
	}

	manualES := NewXLanguage("messages", language.Spanish)
	for id, val := range teststrings {
		manualES.Set(id, val)
	}

	// Load from xml file
	loadxmlES, err := NewXLanguageFromXMLFile("./testunit/errors.es.xml")
	if err != nil {
		t.Errorf("Error loading XML File into language %s", err)
		return
	}

	// Load from text file
	loadtextES, err := NewXLanguageFromFile("./testunit/errors.es.txt")
	if err != nil {
		t.Errorf("Error loading Text File into language %s", err)
		return
	}

	return

	// from xml string
	xmlstr, err := ioutil.ReadFile("./testunit/errors.es.xml")
	if err != nil {
		t.Errorf("Error loading XML File errors into string %s", err)
		return
	}
	loadstringxmlES, err := NewXLanguageFromXMLString(string(xmlstr))
	if err != nil {
		t.Errorf("Error loading XML File into language %s", err)
		return
	}

	// from text string
	textstr, err := ioutil.ReadFile("./testunit/errors.es.txt")
	if err != nil {
		t.Errorf("Error loading Text File errors into string %s", err)
		return
	}
	loadstringtxtES, err := NewXLanguageFromString(string(textstr))
	if err != nil {
		t.Errorf("Error loading Text File into language %s", err)
		return
	}

	// verify all
	for id, val := range teststrings {
		v1 := manualES.Get(id)
		if v1 != val {
			t.Errorf("Error reading value of manualES::%s", v1)
			return
		}
		v2 := loadxmlES.Get(id)
		if v2 != val {
			t.Errorf("Error reading value of loadxmlES::%s", v2)
			return
		}
		v3 := loadtextES.Get(id)
		if v3 != val {
			t.Errorf("Error reading value of loadtextES::%s", v3)
			return
		}
		v4 := loadstringxmlES.Get(id)
		if v4 != val {
			t.Errorf("Error reading value of loadstringxmlES::%s", v4)
			return
		}
		v5 := loadstringtxtES.Get(id)
		if v5 != val {
			t.Errorf("Error reading value of loadstringtxtES::%s", v5)
			return
		}
	}
}

func TestXLanguageAssign(t *testing.T) {

	// Those strings are the same in errors.es.txt and errors.es.xml to make the tests
	teststrings := map[string]string{
		"panicerror":  "System critical error",
		"systemerror": "System error",
		"fileerror":   "File error",
	}

	manualES := NewXLanguage("messages", language.Spanish)
	for id, val := range teststrings {
		manualES.Set(id, val)
	}

	manualES.SetName("errors")
	manualES.SetLanguage(language.Spanish)
	manualES.Set("noerror", "There is no error")

	// Overload the spanish entries
	err := manualES.LoadFile("./testunit/errors.es.txt")
	if err != nil {
		t.Errorf("Error loading Text File into language %s", err)
		return
	}

	v1 := manualES.Get("panicerror")
	if v1 != "Error crítico del sistema" {
		t.Errorf("The value of panicerror is not correct %s", v1)
		return
	}

	v2 := manualES.Get("noerror") // should still be in english
	if v2 != "There is no error" {
		t.Errorf("The value of noerror is not correct %s", v2)
		return
	}

	manualES.Del("noerror")
	v3 := manualES.Get("noerror") // should be empty
	if v3 != "" {
		t.Errorf("The value of noerror is not correct (empty) %s", v3)
		return
	}

	name := manualES.GetName()
	if name != "errors" {
		t.Errorf("The value of language name is not correct %s", name)
		return
	}

	lang := manualES.GetLanguage()
	if lang != language.Spanish {
		t.Errorf("The value of language is not correct %s", lang)
		return
	}

	// String
	str := fmt.Sprint(manualES)
	if str != "xcore.XLanguage{fileerror:Error leyendo el archivo panicerror:Error crítico del sistema systemerror:Error del sistema}" {
		t.Errorf("The print value language is not correct %s", str)
		return
	}

	// Gostring
	str = fmt.Sprintf("%#v", manualES)
	if str != "#xcore.XLanguage{fileerror:\"Error leyendo el archivo\" panicerror:\"Error crítico del sistema\" systemerror:\"Error del sistema\"}" {
		t.Errorf("The print #value language is not correct %s", str)
		return
	}
}
