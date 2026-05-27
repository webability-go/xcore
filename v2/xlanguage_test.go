package xcore

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"golang.org/x/text/language"
)

// ---------------------------------------------------------------------------
// Example functions (appear in go.dev documentation)
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// Core functionality tests
// ---------------------------------------------------------------------------

// TestXLanguageBasic verifies that all four loading paths (XML file, text file,
// XML string, text string) produce the expected entries.
func TestXLanguageBasic(t *testing.T) {

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
		t.Fatalf("Error loading XML File into language: %s", err)
	}

	// Load from text file
	loadtextES, err := NewXLanguageFromFile("./testunit/errors.es.txt")
	if err != nil {
		t.Fatalf("Error loading Text File into language: %s", err)
	}

	// from xml string
	xmlbytes, err := os.ReadFile("./testunit/errors.es.xml")
	if err != nil {
		t.Fatalf("Error reading XML file into string: %s", err)
	}
	loadstringxmlES, err := NewXLanguageFromXMLString(string(xmlbytes))
	if err != nil {
		t.Fatalf("Error loading XML string into language: %s", err)
	}

	// from text string
	txtbytes, err := os.ReadFile("./testunit/errors.es.txt")
	if err != nil {
		t.Fatalf("Error reading text file into string: %s", err)
	}
	loadstringtxtES, err := NewXLanguageFromString(string(txtbytes))
	if err != nil {
		t.Fatalf("Error loading text string into language: %s", err)
	}

	for id, expected := range teststrings {
		if v := manualES.Get(id); v != expected {
			t.Errorf("manualES.Get(%q) = %q, want %q", id, v, expected)
		}
		if v := loadxmlES.Get(id); v != expected {
			t.Errorf("loadxmlES.Get(%q) = %q, want %q", id, v, expected)
		}
		if v := loadtextES.Get(id); v != expected {
			t.Errorf("loadtextES.Get(%q) = %q, want %q", id, v, expected)
		}
		if v := loadstringxmlES.Get(id); v != expected {
			t.Errorf("loadstringxmlES.Get(%q) = %q, want %q", id, v, expected)
		}
		if v := loadstringtxtES.Get(id); v != expected {
			t.Errorf("loadstringtxtES.Get(%q) = %q, want %q", id, v, expected)
		}
	}
}

// TestXLanguageAssign tests SetName / SetLanguage / GetName / GetLanguage,
// overlaying loads, Del, String, and GoString.
func TestXLanguageAssign(t *testing.T) {

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

	// Overload the Spanish entries with real translations
	if err := manualES.LoadFile("./testunit/errors.es.txt"); err != nil {
		t.Fatalf("Error loading Text File into language: %s", err)
	}

	if v := manualES.Get("panicerror"); v != "Error crítico del sistema" {
		t.Errorf("panicerror = %q, want 'Error crítico del sistema'", v)
	}

	// noerror must survive the overlay
	if v := manualES.Get("noerror"); v != "There is no error" {
		t.Errorf("noerror = %q, want 'There is no error'", v)
	}

	manualES.Del("noerror")
	if v := manualES.Get("noerror"); v != "" {
		t.Errorf("after Del, noerror = %q, want ''", v)
	}

	if name := manualES.GetName(); name != "errors" {
		t.Errorf("GetName() = %q, want 'errors'", name)
	}

	if lang := manualES.GetLanguage(); lang != language.Spanish {
		t.Errorf("GetLanguage() = %v, want Spanish", lang)
	}

	// String() output
	wantStr := "xcore.XLanguage{fileerror:Error leyendo el archivo panicerror:Error crítico del sistema systemerror:Error del sistema}"
	if got := fmt.Sprint(manualES); got != wantStr {
		t.Errorf("String() = %q, want %q", got, wantStr)
	}

	// GoString() output
	wantGoStr := "#xcore.XLanguage{fileerror:\"Error leyendo el archivo\" panicerror:\"Error crítico del sistema\" systemerror:\"Error del sistema\"}"
	if got := fmt.Sprintf("%#v", manualES); got != wantGoStr {
		t.Errorf("GoString() = %q, want %q", got, wantGoStr)
	}
}

// TestXLanguageGetXMLRoundTrip verifies that GetXML() produces valid XML that
// can be reloaded by LoadXMLString with identical results (fix 2.3).
func TestXLanguageGetXMLRoundTrip(t *testing.T) {
	orig, err := NewXLanguageFromXMLFile("./testunit/errors.es.xml")
	if err != nil {
		t.Fatalf("NewXLanguageFromXMLFile: %v", err)
	}

	xmlOut := orig.GetXML()
	if xmlOut == "" {
		t.Fatal("GetXML() returned empty string")
	}

	reloaded, err := NewXLanguageFromXMLString(xmlOut)
	if err != nil {
		t.Fatalf("NewXLanguageFromXMLString on GetXML output: %v\nXML was:\n%s", err, xmlOut)
	}

	for id, val := range orig.GetEntries() {
		if got := reloaded.Get(id); got != val {
			t.Errorf("round-trip mismatch for %q: got %q, want %q", id, got, val)
		}
	}
}

// TestXLanguageGetXMLSpecialChars verifies that keys, status values, and entry
// values containing XML special characters are correctly escaped and survive a
// round-trip (fix 2.3 – injection prevention).
func TestXLanguageGetXMLSpecialChars(t *testing.T) {
	l := NewXLanguage("test&name", language.English)
	l.Set(`key"with"quotes`, `value <with> & "special" chars`)
	l.Set(`key<evil>`, `another & "test"`)
	l.SetStatus(`key"with"quotes`, `status<>&"`)

	xmlOut := l.GetXML()

	reloaded, err := NewXLanguageFromXMLString(xmlOut)
	if err != nil {
		t.Fatalf("special-chars round-trip: NewXLanguageFromXMLString failed: %v\nXML:\n%s", err, xmlOut)
	}

	cases := map[string]string{
		`key"with"quotes`: `value <with> & "special" chars`,
		`key<evil>`:       `another & "test"`,
	}
	for key, want := range cases {
		if got := reloaded.Get(key); got != want {
			t.Errorf("special char round-trip: Get(%q) = %q, want %q", key, got, want)
		}
	}

	if got := reloaded.GetStatus(`key"with"quotes`); got != `status<>&"` {
		t.Errorf("special char round-trip: GetStatus = %q, want 'status<>&\"'", got)
	}
}

// TestXLanguageGetXMLCDATATerminator verifies that a value containing ]]>
// (which would break a raw CDATA block) still survives a round-trip (fix 2.3).
func TestXLanguageGetXMLCDATATerminator(t *testing.T) {
	l := NewXLanguage("test", language.English)
	tricky := "before]]>after"
	l.Set("tricky", tricky)

	xmlOut := l.GetXML()
	reloaded, err := NewXLanguageFromXMLString(xmlOut)
	if err != nil {
		t.Fatalf("CDATA terminator round-trip failed: %v\nXML:\n%s", err, xmlOut)
	}
	if got := reloaded.Get("tricky"); got != tricky {
		t.Errorf("CDATA terminator: Get(tricky) = %q, want %q", got, tricky)
	}
}

// TestXLanguageInvalidXML verifies that malformed XML is correctly rejected
// (fix 2.6 – strict decoder).
func TestXLanguageInvalidXML(t *testing.T) {
	bad := []string{
		`not xml at all`,
		`<language id="x" lang="en"><entry id="k">unclosed`,
		`<language id="x" lang="en"><entry id="k">val</entry>`,           // missing closing </language>
		`<language id="x" lang="en"><entry id="k">ok</ENTRY></language>`, // mismatched tag
	}
	for _, input := range bad {
		_, err := NewXLanguageFromXMLString(input)
		if err == nil {
			t.Errorf("expected error for invalid XML %q, got nil", input)
		}
	}
}

// TestXLanguageInvalidLanguageTag verifies that an unrecognised language tag
// does not panic and falls back to language.Und (fix 2.2).
func TestXLanguageInvalidLanguageTag(t *testing.T) {
	xmldata := `<?xml version="1.0" encoding="UTF-8"?>
<language id="test" lang="zz-INVALID-TAG">
  <entry id="hello">world</entry>
</language>`

	l, err := NewXLanguageFromXMLString(xmldata)
	if err != nil {
		t.Fatalf("unexpected error with invalid language tag: %v", err)
	}
	// Should fall back to language.Und and still load entries
	if got := l.Get("hello"); got != "world" {
		t.Errorf("Get(hello) = %q, want 'world'", got)
	}
	if got := l.GetLanguage(); got != language.Und {
		t.Errorf("GetLanguage() = %v, want Und", got)
	}
}

// TestXLanguageSetGetNameLanguage verifies that SetName/GetName and
// SetLanguage/GetLanguage are consistent and protected (fix 2.1).
func TestXLanguageSetGetNameLanguage(t *testing.T) {
	l := NewXLanguage("initial", language.English)

	if got := l.GetName(); got != "initial" {
		t.Errorf("GetName() = %q, want 'initial'", got)
	}
	if got := l.GetLanguage(); got != language.English {
		t.Errorf("GetLanguage() = %v, want English", got)
	}

	l.SetName("updated")
	l.SetLanguage(language.French)

	if got := l.GetName(); got != "updated" {
		t.Errorf("after SetName, GetName() = %q, want 'updated'", got)
	}
	if got := l.GetLanguage(); got != language.French {
		t.Errorf("after SetLanguage, GetLanguage() = %v, want French", got)
	}
}

// TestXLanguageStatus verifies SetStatus / GetStatus.
func TestXLanguageStatus(t *testing.T) {
	l := NewXLanguage("test", language.English)
	l.Set("key1", "value1")
	l.SetStatus("key1", "approved")

	if got := l.GetStatus("key1"); got != "approved" {
		t.Errorf("GetStatus(key1) = %q, want 'approved'", got)
	}
	if got := l.GetStatus("missing"); got != "" {
		t.Errorf("GetStatus(missing) = %q, want ''", got)
	}
}

// TestXLanguageGetEntries verifies that GetEntries returns a full copy and
// mutations of the copy do not affect the original.
func TestXLanguageGetEntries(t *testing.T) {
	l := NewXLanguage("test", language.English)
	l.Set("a", "1")
	l.Set("b", "2")

	entries := l.GetEntries()
	if len(entries) != 2 {
		t.Fatalf("GetEntries() len = %d, want 2", len(entries))
	}

	// Mutate the copy; original must be unaffected
	entries["a"] = "MUTATED"
	if got := l.Get("a"); got != "1" {
		t.Errorf("after mutating copy, Get(a) = %q, want '1'", got)
	}
}

// TestXLanguageLoadStringBatchAtomicity verifies that after LoadString
// all entries are visible together (fix 2.4 – single-lock batch update).
func TestXLanguageLoadStringBatchAtomicity(t *testing.T) {
	const n = 100
	lines := make([]string, n)
	for i := 0; i < n; i++ {
		lines[i] = fmt.Sprintf("key%d=value%d", i, i)
	}
	data := strings.Join(lines, "\n")

	l := NewXLanguage("test", language.English)
	if err := l.LoadString(data); err != nil {
		t.Fatalf("LoadString: %v", err)
	}
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("key%d", i)
		want := fmt.Sprintf("value%d", i)
		if got := l.Get(key); got != want {
			t.Errorf("Get(%q) = %q, want %q", key, got, want)
		}
	}
}

// TestXLanguageFileSizeLimit verifies that files exceeding maxLanguageFileSize
// are rejected with an error (fix 2.5).
func TestXLanguageFileSizeLimit(t *testing.T) {
	// Write a temporary file larger than maxLanguageFileSize
	tmp, err := os.CreateTemp(t.TempDir(), "xlang_big_*.txt")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	defer tmp.Close()

	// Write slightly more than the limit
	chunk := make([]byte, 1024*1024) // 1 MB
	for i := range chunk {
		chunk[i] = 'a'
	}
	for i := 0; i < 11; i++ { // 11 MB total → exceeds 10 MB limit
		if _, err := tmp.Write(chunk); err != nil {
			t.Fatalf("writing temp file: %v", err)
		}
	}
	tmp.Close()

	l := NewXLanguage("test", language.English)
	if err := l.LoadFile(tmp.Name()); err == nil {
		t.Error("LoadFile on oversized file should return an error, got nil")
	}
	if err := l.LoadXMLFile(tmp.Name()); err == nil {
		t.Error("LoadXMLFile on oversized file should return an error, got nil")
	}
}

// TestXLanguageStringSizeLimit verifies that strings exceeding maxLanguageFileSize
// are rejected by LoadXMLString (fix 2.5).
func TestXLanguageStringSizeLimit(t *testing.T) {
	huge := strings.Repeat("a", maxLanguageFileSize+1)
	l := NewXLanguage("test", language.English)
	if err := l.LoadXMLString(huge); err == nil {
		t.Error("LoadXMLString with oversized string should return an error, got nil")
	}
}

// TestXLanguageMissingFile verifies that loading from a non-existent file
// returns an error and does not panic.
func TestXLanguageMissingFile(t *testing.T) {
	_, err := NewXLanguageFromXMLFile("/nonexistent/path/file.xml")
	if err == nil {
		t.Error("expected error for missing XML file, got nil")
	}
	_, err = NewXLanguageFromFile("/nonexistent/path/file.txt")
	if err == nil {
		t.Error("expected error for missing text file, got nil")
	}
}

// TestXLanguageLoadStringComments verifies that comment and blank lines are
// skipped in the key=value format.
func TestXLanguageLoadStringComments(t *testing.T) {
	data := `
# this is a comment
; this too
valid=hello

another=world
=ignored
`
	l, err := NewXLanguageFromString(data)
	if err != nil {
		t.Fatalf("NewXLanguageFromString: %v", err)
	}

	if got := l.Get("valid"); got != "hello" {
		t.Errorf("Get(valid) = %q, want 'hello'", got)
	}
	if got := l.Get("another"); got != "world" {
		t.Errorf("Get(another) = %q, want 'world'", got)
	}
	// empty key must be ignored
	if got := l.Get(""); got != "" {
		t.Errorf("Get('') = %q, want ''", got)
	}
}

// ---------------------------------------------------------------------------
// Concurrency tests (fix 2.1 – mutex coverage for Name, Language and entries)
// ---------------------------------------------------------------------------

// TestXLanguageConcurrentSetGet fires many goroutines doing concurrent
// Set / Get operations and verifies there is no data race.
// Run with: go test -race
func TestXLanguageConcurrentSetGet(t *testing.T) {
	l := NewXLanguage("concurrent", language.English)
	const goroutines = 50
	const iterations = 200

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				key := fmt.Sprintf("key%d", id%10)
				l.Set(key, fmt.Sprintf("val%d-%d", id, j))
			}
		}(i)

		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = l.Get(fmt.Sprintf("key%d", id%10))
			}
		}(i)
	}
	wg.Wait()
}

// TestXLanguageConcurrentLoadXMLString races multiple LoadXMLString calls to
// verify that Name and Language mutations are protected by the mutex (fix 2.1).
func TestXLanguageConcurrentLoadXMLString(t *testing.T) {
	xmlES := `<?xml version="1.0" encoding="UTF-8"?>
<language id="errors_es" lang="es">
  <entry id="key1">valor_es</entry>
</language>`

	xmlEN := `<?xml version="1.0" encoding="UTF-8"?>
<language id="errors_en" lang="en">
  <entry id="key1">value_en</entry>
</language>`

	l := NewXLanguage("initial", language.Und)
	const goroutines = 30

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			src := xmlES
			if i%2 == 0 {
				src = xmlEN
			}
			if err := l.LoadXMLString(src); err != nil {
				t.Errorf("concurrent LoadXMLString: %v", err)
			}
			// Just read to exercise the RLock path
			_ = l.GetName()
			_ = l.GetLanguage()
			_ = l.Get("key1")
		}(i)
	}
	wg.Wait()
}

// TestXLanguageConcurrentSetNameLanguage races SetName/GetName and
// SetLanguage/GetLanguage to confirm mutex protection (fix 2.1).
func TestXLanguageConcurrentSetNameLanguage(t *testing.T) {
	l := NewXLanguage("start", language.English)
	const goroutines = 40

	var wg sync.WaitGroup
	wg.Add(goroutines)

	langs := []language.Tag{language.English, language.French, language.Spanish, language.German}

	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				l.SetName(fmt.Sprintf("name%d", i))
				l.SetLanguage(langs[i%len(langs)])
			} else {
				_ = l.GetName()
				_ = l.GetLanguage()
			}
		}(i)
	}
	wg.Wait()
}

// TestXLanguageConcurrentLoadString races many LoadString calls to verify that
// the single-lock batch update does not cause intermediate visibility (fix 2.4).
func TestXLanguageConcurrentLoadString(t *testing.T) {
	l := NewXLanguage("test", language.English)
	const goroutines = 20

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			data := fmt.Sprintf("gkey%d=gval%d\n", i, i)
			if err := l.LoadString(data); err != nil {
				t.Errorf("concurrent LoadString: %v", err)
			}
		}(i)
	}
	wg.Wait()
}

// ---------------------------------------------------------------------------
// XML round-trip with Name preserved in GetXML (fix 2.3)
// ---------------------------------------------------------------------------

// TestXLanguageGetXMLNamePreserved checks that the language name and tag
// survive the GetXML → LoadXMLString round-trip.
func TestXLanguageGetXMLNamePreserved(t *testing.T) {
	l := NewXLanguage("my-language", language.French)
	l.Set("greeting", "Bonjour")

	xmlOut := l.GetXML()
	reloaded, err := NewXLanguageFromXMLString(xmlOut)
	if err != nil {
		t.Fatalf("NewXLanguageFromXMLString: %v", err)
	}

	if got := reloaded.GetName(); got != "my-language" {
		t.Errorf("GetName() after round-trip = %q, want 'my-language'", got)
	}
	if got := reloaded.GetLanguage(); got != language.French {
		t.Errorf("GetLanguage() after round-trip = %v, want French", got)
	}
	if got := reloaded.Get("greeting"); got != "Bonjour" {
		t.Errorf("Get(greeting) after round-trip = %q, want 'Bonjour'", got)
	}
}

// ---------------------------------------------------------------------------
// File-backed helpers
// ---------------------------------------------------------------------------

// TestXLanguageFileRoundTrip writes an XML file, reads it back, and checks that
// all entries match (exercises LoadXMLFile with os.ReadFile path, fix 2.7).
func TestXLanguageFileRoundTrip(t *testing.T) {
	l := NewXLanguage("test", language.English)
	l.Set("hello", "world")
	l.Set("foo", "bar")

	dir := t.TempDir()
	path := filepath.Join(dir, "lang.xml")
	if err := os.WriteFile(path, []byte(l.GetXML()), 0o644); err != nil {
		t.Fatalf("writing temp XML file: %v", err)
	}

	loaded, err := NewXLanguageFromXMLFile(path)
	if err != nil {
		t.Fatalf("NewXLanguageFromXMLFile: %v", err)
	}

	for id, want := range map[string]string{"hello": "world", "foo": "bar"} {
		if got := loaded.Get(id); got != want {
			t.Errorf("file round-trip: Get(%q) = %q, want %q", id, got, want)
		}
	}
}

// TestXLanguageEmptyGet verifies that Get returns "" for unknown keys.
func TestXLanguageEmptyGet(t *testing.T) {
	l := NewXLanguage("test", language.English)
	if got := l.Get("nonexistent"); got != "" {
		t.Errorf("Get(nonexistent) = %q, want ''", got)
	}
}
