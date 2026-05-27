package xcore

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/text/language"
)

// ---------------------------------------------------------------------------
// Example functions (appear in go.dev documentation)
// ---------------------------------------------------------------------------

func ExampleNewXTemplateFromString() {
	tmpl, _ := NewXTemplateFromString(`
%-- This is a comment. It will not appear in the final code. --%
Let's put your name here: {{clientname}}<br />
And lets put your hobbies here:<br />
%-- note the 1rst id is the entry into the data to inject and the second one is the name of the sub-template to use --%
@@hobbies:hobby@@
%-- And you need the template for each hobby:--%
[[hobby]]
I love {{name}}<br />
[[]]
`)
	data := XDataset{
		"clientname": "Fred",
		"hobbies": &XDatasetCollection{
			&XDataset{"name": "Football"},
			&XDataset{"name": "Ping-pong"},
			&XDataset{"name": "Swimming"},
			&XDataset{"name": "Videogames"},
		},
	}
	fmt.Println(tmpl.Execute(&data))
	// Output:
	// Let's put your name here: Fred<br />
	// And lets put your hobbies here:<br />
	// I love Football<br />
	// I love Ping-pong<br />
	// I love Swimming<br />
	// I love Videogames<br />
}

// ---------------------------------------------------------------------------
// Core metalanguage: {{variable}}
// ---------------------------------------------------------------------------

func TestMetaVariable(t *testing.T) {
	tmpl, err := NewXTemplateFromString("Hello {{name}}!")
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	data := XDataset{"name": "World"}
	if got := tmpl.Execute(&data); got != "Hello World!" {
		t.Errorf("got %q", got)
	}
}

func TestMetaVariableMissing(t *testing.T) {
	tmpl, _ := NewXTemplateFromString("Hello {{missing}}!")
	data := XDataset{}
	if got := tmpl.Execute(&data); got != "Hello !" {
		t.Errorf("missing var should be empty, got %q", got)
	}
}

func TestMetaVariableNestedPath(t *testing.T) {
	tmpl, _ := NewXTemplateFromString("{{user>name}}")
	data := XDataset{
		"user": &XDataset{"name": "Alice"},
	}
	if got := tmpl.Execute(&data); got != "Alice" {
		t.Errorf("nested path: got %q", got)
	}
}

func TestMetaVariableNilDatacol(t *testing.T) {
	tmpl, _ := NewXTemplateFromString("{{name}}")
	// Execute with nil data goes through injectorDepth(nil,nil,0)
	if got := tmpl.Execute(nil); got != "" {
		t.Errorf("nil datacol: got %q, want ''", got)
	}
}

// ---------------------------------------------------------------------------
// Core metalanguage: %-- comment --%
// ---------------------------------------------------------------------------

func TestMetaComment(t *testing.T) {
	tmpl1, _ := NewXTemplateFromString("abcdefg")
	tmpl2, _ := NewXTemplateFromString(`a%--comment1--%b%--comment 2
[[]]
[[subtemplate]]
		--%c%--comment3 --%
defg%-- ending @@subtemplate@@ comment --%
`)
	if tmpl1.Execute(nil) != tmpl2.Execute(nil) {
		t.Errorf("comments not stripped correctly:\ntmpl1=%q\ntmpl2=%q",
			tmpl1.Execute(nil), tmpl2.Execute(nil))
	}
}

// ---------------------------------------------------------------------------
// Core metalanguage: ##language##
// ---------------------------------------------------------------------------

func TestMetaLanguage(t *testing.T) {
	tmpl, _ := NewXTemplateFromString("Test with ##some## ##languages## here")
	data := &XDataset{}
	l, _ := NewXLanguageFromString("some=a tiny table\nlanguages=of english language\n")
	data.Set("#", l)
	if got := tmpl.Execute(data); got != "Test with a tiny table of english language here" {
		t.Errorf("language injection: got %q", got)
	}
}

func TestMetaLanguageMissing(t *testing.T) {
	tmpl, _ := NewXTemplateFromString("##missing##")
	data := &XDataset{}
	l, _ := NewXLanguageFromString("other=value\n")
	data.Set("#", l)
	// missing language key -> empty string
	if got := tmpl.Execute(data); got != "" {
		t.Errorf("missing language key: got %q, want ''", got)
	}
}

func TestMetaLanguageNil(t *testing.T) {
	tmpl, _ := NewXTemplateFromString("##key##")
	data := &XDataset{} // no "#" entry
	// should not panic, just produce empty string
	if got := tmpl.Execute(data); got != "" {
		t.Errorf("nil language: got %q, want ''", got)
	}
}

// ---------------------------------------------------------------------------
// Core metalanguage: &&reference&&
// ---------------------------------------------------------------------------

func TestMetaReference(t *testing.T) {
	tmpl, err := NewXTemplateFromString(`
The sub template starts here: &&template1&&. End.
[[template1]]
This is the template 1
[[]]
`)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	got := tmpl.Execute(&XDataset{})
	if !strings.Contains(got, "This is the template 1") {
		t.Errorf("reference not injected: %q", got)
	}
}

func TestMetaReferenceWithDataField(t *testing.T) {
	// &&name:field&& — pushes datacol.GetData(field) as new layer
	tmpl, _ := NewXTemplateFromString(`&&info:meta&&
[[info]]{{preferred-color}}[[]]`)
	data := XDataset{
		"meta": &XDataset{"preferred-color": "blue"},
	}
	// The \n between &&info:meta&& and [[info]] is outer-template text and is rendered.
	if got := tmpl.Execute(&data); got != "blue\n" {
		t.Errorf("reference with data field: got %q, want 'blue\\n'", got)
	}
}

func TestMetaReferenceDynamic(t *testing.T) {
	// &&:field:prefix&& resolves template name from data value
	tmpl, _ := NewXTemplateFromString(`??sport:sport??
[[sport.yes]]It is a sport![[]]
[[sport.no]]Not a sport.[[]]
[[sport]]Unknown.[[]]`)
	// The \n between ??sport:sport?? and [[sport.yes]] is outer-template text.
	for _, tc := range []struct{ sport, want string }{
		{"yes", "It is a sport!\n"},
		{"no", "Not a sport.\n"},
		{"maybe", "Unknown.\n"},
	} {
		data := XDataset{"sport": tc.sport}
		if got := tmpl.Execute(&data); got != tc.want {
			t.Errorf("sport=%q: got %q, want %q", tc.sport, got, tc.want)
		}
	}
}

func TestMetaReferenceC(t *testing.T) {
	tmpl, err := NewXTemplateFromFile("testunit/c.template")
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	data := &XDataset{}
	if got := tmpl.Execute(data); got != "ABC\nDA\nEAF\n\n" {
		t.Errorf("c.template reference: got %q", got)
	}
}

// ---------------------------------------------------------------------------
// Core metalanguage: @@range@@
// ---------------------------------------------------------------------------

func TestMetaRangeBasic(t *testing.T) {
	tmpl, _ := NewXTemplateFromString(`@@hobbies:hobby@@
[[hobby]]
I love {{name}}<br />
[[]]
`)
	data := XDataset{
		"hobbies": &XDatasetCollection{
			&XDataset{"name": "Football"},
			&XDataset{"name": "Ping-pong"},
		},
	}
	got := tmpl.Execute(&data)
	if !strings.Contains(got, "I love Football") || !strings.Contains(got, "I love Ping-pong") {
		t.Errorf("range basic: %q", got)
	}
}

func TestMetaRangeCounter(t *testing.T) {
	tmpl, _ := NewXTemplateFromString(`@@items:item@@
[[item]]{{.counter}}.{{name}} [[]]`)
	data := XDataset{
		"items": &XDatasetCollection{
			&XDataset{"name": "a"},
			&XDataset{"name": "b"},
			&XDataset{"name": "c"},
		},
	}
	// The \n between @@items:item@@ and [[item]] is outer-template text rendered after the range.
	if got := tmpl.Execute(&data); got != "1.a 2.b 3.c \n" {
		t.Errorf("range counter: got %q", got)
	}
}

func TestMetaRangeFirstLastEven(t *testing.T) {
	tmpl, _ := NewXTemplateFromString(`@@items:item@@
[[item.first]]F{{name}} [[]]
[[item.last]]L{{name}} [[]]
[[item.even]]E{{name}} [[]]
[[item]]{{name}} [[]]
`)
	data := XDataset{
		"items": &XDatasetCollection{
			&XDataset{"name": "a"}, // first
			&XDataset{"name": "b"}, // even (index 1 is odd, index 0 even... wait, i%2==0 is even so 0,2,4)
			&XDataset{"name": "c"}, // last
		},
	}
	got := tmpl.Execute(&data)
	// index 0 -> .first, index 1 -> .even (i%2==0? no i=1, 1%2=1 so not even) -> plain, index 2 -> .last
	if !strings.Contains(got, "Fa") {
		t.Errorf("range first: missing 'Fa' in %q", got)
	}
	if !strings.Contains(got, "Lc") {
		t.Errorf("range last: missing 'Lc' in %q", got)
	}
}

func TestMetaRangeKeyN(t *testing.T) {
	tmpl, _ := NewXTemplateFromString(`@@items:item@@
[[item.key.1]]KEY1:{{name}} [[]]
[[item]]{{name}} [[]]
`)
	data := XDataset{
		"items": &XDatasetCollection{
			&XDataset{"name": "a"},
			&XDataset{"name": "b"}, // key.1
			&XDataset{"name": "c"},
		},
	}
	got := tmpl.Execute(&data)
	if !strings.Contains(got, "KEY1:b") {
		t.Errorf("range key.N: missing 'KEY1:b' in %q", got)
	}
}

func TestMetaRangeEmpty(t *testing.T) {
	tmpl, _ := NewXTemplateFromString(`@@hobbies@@
[[hobbies]]I love {{name}}<br />[[]]
[[hobbies.none]]I love nothing<br />[[]]
`)
	data := XDataset{"hobbies": &XDatasetCollection{}}
	got := tmpl.Execute(&data)
	if !strings.Contains(got, "I love nothing") {
		t.Errorf("range empty .none: %q", got)
	}
}

func TestMetaRangeNilCollection(t *testing.T) {
	tmpl, _ := NewXTemplateFromString(`@@hobbies@@
[[hobbies]]{{name}}[[]]
[[hobbies.none]]nothing[[]]
`)
	var col XDatasetCollection
	data := XDataset{"hobbies": col}
	got := tmpl.Execute(&data)
	if !strings.Contains(got, "nothing") {
		t.Errorf("range nil collection .none: %q", got)
	}
}

// ---------------------------------------------------------------------------
// Core metalanguage: ??condition??
// ---------------------------------------------------------------------------

func TestMetaConditionTrue(t *testing.T) {
	tmpl, _ := NewXTemplateFromString(`??name??
[[name]]Hello {{name}}![[]]`)
	// The \n between ??name?? and [[name]] is outer-template text rendered after the condition.
	data := XDataset{"name": "Alice"}
	if got := tmpl.Execute(&data); got != "Hello Alice!\n" {
		t.Errorf("condition true: got %q", got)
	}
}

func TestMetaConditionFalse(t *testing.T) {
	tmpl, _ := NewXTemplateFromString(`??name??
[[name]]Hello {{name}}![[]]
[[name.none]]nobody here[[]]`)
	// The \n between ??name?? and [[name]] is outer-template text rendered after the condition.
	data := XDataset{} // name absent
	got := tmpl.Execute(&data)
	if got != "nobody here\n" {
		t.Errorf("condition false (.none): got %q", got)
	}
}

func TestMetaConditionNoneTemplateFix(t *testing.T) {
	// Regression test for the original bug where .none template was found
	// but subt was used instead of tmp. Verify .none IS executed when present.
	tmpl, _ := NewXTemplateFromString(`TAGADA
??familyname??
[[familyname]]
Hello {{familyname}}<br />
[[]]
[[familyname.none]]
I don't know you<br />
[[]]
`)
	data := XDataset{"clientname": "Fred"} // no familyname
	got := tmpl.Execute(&data)
	if !strings.Contains(got, "I don't know you") {
		t.Errorf(".none template not used: got %q", got)
	}
	if strings.Contains(got, "Hello") {
		t.Errorf("wrong template used: got %q", got)
	}
}

func TestMetaConditionValueSpecific(t *testing.T) {
	// ??field?? picks subtemplate named field.value when present
	tmpl, _ := NewXTemplateFromString(`??status??
[[status.active]]ACTIVE[[]]
[[status.inactive]]INACTIVE[[]]
[[status]]OTHER:{{status}}[[]]`)
	// The \n between ??status?? and [[status.active]] is outer-template text rendered after the condition.
	for _, tc := range []struct{ status, want string }{
		{"active", "ACTIVE\n"},
		{"inactive", "INACTIVE\n"},
		{"pending", "OTHER:pending\n"},
	} {
		data := XDataset{"status": tc.status}
		if got := tmpl.Execute(&data); got != tc.want {
			t.Errorf("status=%q: got %q, want %q", tc.status, got, tc.want)
		}
	}
}

func TestMetaConditionWithDataset(t *testing.T) {
	tmpl, _ := NewXTemplateFromString(`??meta??
[[meta]]color={{preferred-color}}[[]]`)
	data := XDataset{
		"meta": &XDataset{"preferred-color": "green"},
	}
	// The \n between ??meta?? and [[meta]] is outer-template text rendered after the condition.
	if got := tmpl.Execute(&data); got != "color=green\n" {
		t.Errorf("condition with dataset: got %q", got)
	}
}

// ---------------------------------------------------------------------------
// Core metalanguage: !!dump!!
// ---------------------------------------------------------------------------

func TestMetaDump(t *testing.T) {
	tmpl, _ := NewXTemplateFromString("!!dump!!")
	data := XDataset{"key": "val"}
	got := tmpl.Execute(&data)
	// Should contain the GoString of the dataset
	if got == "" {
		t.Errorf("dump: got empty string")
	}
}

func TestMetaDumpNilDatacol(t *testing.T) {
	tmpl, _ := NewXTemplateFromString("!!dump!!")
	// Execute with nil -> injectorDepth(nil,nil,0) -> datacol nil -> nothing output
	if got := tmpl.Execute(nil); got != "" {
		t.Errorf("dump with nil datacol: got %q", got)
	}
}

// ---------------------------------------------------------------------------
// Output filters (fix 3.1)
// ---------------------------------------------------------------------------

func TestFilterHTML(t *testing.T) {
	tmpl, _ := NewXTemplateFromString("{{name|html}}")
	cases := []struct{ input, want string }{
		{`<script>alert("xss")</script>`, `&lt;script&gt;alert(&#34;xss&#34;)&lt;/script&gt;`},
		{`a & b`, `a &amp; b`},
		{`plain`, `plain`},
	}
	for _, tc := range cases {
		data := XDataset{"name": tc.input}
		if got := tmpl.Execute(&data); got != tc.want {
			t.Errorf("html filter %q: got %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestFilterJS(t *testing.T) {
	tmpl, _ := NewXTemplateFromString("{{val|js}}")
	cases := []struct{ input, want string }{
		{`say "hello"`, `say \"hello\"`},
		{`it's`, `it\'s`},
		{"line\nnewline", `line\nnewline`},
		{"tab\there", `tab\there`},
		{`<tag>`, `\x3Ctag\x3E`},
		{`a & b`, `a \x26 b`},
		{`back\slash`, `back\\slash`},
	}
	for _, tc := range cases {
		data := XDataset{"val": tc.input}
		if got := tmpl.Execute(&data); got != tc.want {
			t.Errorf("js filter %q: got %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestFilterURL(t *testing.T) {
	tmpl, _ := NewXTemplateFromString("{{q|url}}")
	cases := []struct{ input, want string }{
		{"hello world", "hello+world"},
		{"a&b=c", "a%26b%3Dc"},
		{"plain", "plain"},
	}
	for _, tc := range cases {
		data := XDataset{"q": tc.input}
		if got := tmpl.Execute(&data); got != tc.want {
			t.Errorf("url filter %q: got %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestFilterURLQuery(t *testing.T) {
	// urlquery is an alias for url
	tmpl, _ := NewXTemplateFromString("{{q|urlquery}}")
	data := XDataset{"q": "hello world"}
	if got := tmpl.Execute(&data); got != "hello+world" {
		t.Errorf("urlquery filter: got %q", got)
	}
}

func TestFilterRaw(t *testing.T) {
	tmpl, _ := NewXTemplateFromString("{{name|raw}}")
	data := XDataset{"name": `<b>bold</b>`}
	if got := tmpl.Execute(&data); got != `<b>bold</b>` {
		t.Errorf("raw filter: got %q", got)
	}
}

func TestFilterNone(t *testing.T) {
	// No filter = raw pass-through (default)
	tmpl, _ := NewXTemplateFromString("{{name}}")
	data := XDataset{"name": `<b>bold</b>`}
	if got := tmpl.Execute(&data); got != `<b>bold</b>` {
		t.Errorf("no filter: got %q", got)
	}
}

func TestSplitFilter(t *testing.T) {
	cases := []struct{ input, wantPath, wantFilter string }{
		{"name", "name", ""},
		{"name|html", "name", "html"},
		{"a>b>c|js", "a>b>c", "js"},
		{"a|b|url", "a|b", "url"},
		{"|html", "", "html"},
	}
	for _, tc := range cases {
		p, f := splitFilter(tc.input)
		if p != tc.wantPath || f != tc.wantFilter {
			t.Errorf("splitFilter(%q) = (%q,%q), want (%q,%q)",
				tc.input, p, f, tc.wantPath, tc.wantFilter)
		}
	}
}

func TestEscapeJS(t *testing.T) {
	cases := []struct{ input, want string }{
		{`\`, `\\`},
		{`"`, `\"`},
		{`'`, `\'`},
		{"\n", `\n`},
		{"\r", `\r`},
		{"\t", `\t`},
		{"<", `\x3C`},
		{">", `\x3E`},
		{"&", `\x26`},
		{"hello", "hello"},
	}
	for _, tc := range cases {
		if got := escapeJS(tc.input); got != tc.want {
			t.Errorf("escapeJS(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// ExecuteSafe (fix 3.9)
// ---------------------------------------------------------------------------

func TestExecuteSafeCompiledTemplate(t *testing.T) {
	tmpl, _ := NewXTemplateFromString("Hello {{name}}")
	data := XDataset{"name": "World"}
	result, err := tmpl.ExecuteSafe(&data)
	if err != nil {
		t.Fatalf("ExecuteSafe returned unexpected error: %v", err)
	}
	if result != "Hello World" {
		t.Errorf("ExecuteSafe: got %q", result)
	}
}

func TestExecuteSafeNilRoot(t *testing.T) {
	tmpl := &XTemplate{} // Root is nil
	_, err := tmpl.ExecuteSafe(&XDataset{})
	if err == nil {
		t.Error("ExecuteSafe with nil Root should return error")
	}
}

func TestExecuteNilRootReturnsEmpty(t *testing.T) {
	// Execute (backward compat) returns "" when Root is nil, not an error string
	tmpl := &XTemplate{}
	if got := tmpl.Execute(&XDataset{}); got != "" {
		t.Errorf("Execute with nil Root: got %q, want ''", got)
	}
}

// ---------------------------------------------------------------------------
// Recursion depth limit (fix 3.2)
// ---------------------------------------------------------------------------

func TestRecursionDepthLimit(t *testing.T) {
	// Build a chain of 70 mutually-referencing templates — deeper than maxTemplateDepth
	// They should terminate instead of causing a stack overflow.
	root, _ := NewXTemplateFromString("start&&sub0&&end")
	for i := 0; i < 70; i++ {
		name := fmt.Sprintf("sub%d", i)
		next := fmt.Sprintf("sub%d", i+1)
		sub, _ := NewXTemplateFromString(fmt.Sprintf("[%d]&&%s&&", i, next))
		root.AddTemplate(name, sub)
	}
	// sub70 is not defined, so recursion stops there regardless
	// The important thing is it doesn't crash/hang
	got := root.Execute(&XDataset{})
	if !strings.HasPrefix(got, "start") {
		t.Errorf("depth limit: output doesn't start with 'start': %q", got)
	}
}

func TestSelfReferencingTemplate(t *testing.T) {
	// A template that references itself should not infinite-loop
	root, _ := NewXTemplateFromString("&&self&&")
	self, _ := NewXTemplateFromString("x&&self&&")
	root.AddTemplate("self", self)
	self.Father = root // create the cycle: self can find itself via Father
	// Must return within the test timeout (depth limit kicks in)
	done := make(chan struct{})
	go func() {
		root.Execute(&XDataset{})
		close(done)
	}()
	select {
	case <-done:
		// ok
	}
}

// ---------------------------------------------------------------------------
// Clone correctness (fix 3.3)
// ---------------------------------------------------------------------------

func TestCloneIndependence(t *testing.T) {
	orig, _ := NewXTemplateFromString("Hello {{name}}")
	orig.AddTemplate("sub", &XTemplate{
		Root: func() *XTemplateData {
			d := XTemplateData{{ParamType: MetaString, Data: "sub-content"}}
			return &d
		}(),
	})

	cloned := orig.Clone()

	// Executing both should produce the same output
	data := XDataset{"name": "World"}
	if orig.Execute(&data) != cloned.Execute(&data) {
		t.Error("Clone output differs from original")
	}
}

func TestCloneFatherNotCorrupted(t *testing.T) {
	// Fix 3.3: Clone must set Father on the CLONED sub, not the original.
	orig, _ := NewXTemplateFromString("&&sub&&")
	sub, _ := NewXTemplateFromString("original-sub")
	orig.AddTemplate("sub", sub)

	clonedRoot := orig.Clone()

	// The original sub's Father must still be orig, not clonedRoot
	if orig.SubTemplates["sub"].Father != orig {
		t.Errorf("Clone corrupted original sub's Father: got %p, want %p (orig)",
			orig.SubTemplates["sub"].Father, orig)
	}
	// The cloned sub's Father must be clonedRoot
	if clonedRoot.SubTemplates["sub"].Father != clonedRoot {
		t.Errorf("Clone did not set cloned sub's Father to clonedRoot")
	}
}

func TestCloneOutputMatchesOriginal(t *testing.T) {
	orig, err := NewXTemplateFromFile("testunit/b.template")
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	cloned := orig.Clone()

	tmp, _ := time.Parse(time.RFC3339, "2020-01-01T12:00:00")
	lang := NewXLanguage("mainpage", language.English)
	lang.Set("welcome", "Welcome to you")
	data := XDataset{
		"clientname":    "Fred",
		"clientpicture": "face.jpg",
		"hobbies": &XDatasetCollection{
			&XDataset{"name": "Football", "sport": "yes"},
			&XDataset{"name": "Ping-pong", "sport": "yes"},
			&XDataset{"name": "Swimming", "sport": "yes"},
			&XDataset{"name": "Videogames", "sport": "no"},
			&XDataset{"name": "other 5", "sport": nil},
		},
		"preferredhobby": &XDataset{"name": "Baseball", "sport": "yes"},
		"metadata": &XDataset{
			"preferred-color": "blue",
			"Salary":          3568.65,
			"hiredate":        tmp,
		},
		"#": lang,
	}

	s1 := orig.Execute(&data)
	s2 := cloned.Execute(&data)
	if s1 == "" {
		t.Error("original produced empty output")
	}
	if s1 != s2 {
		t.Errorf("clone output differs from original:\noriginal=%q\ncloned=%q", s1, s2)
	}
}

// ---------------------------------------------------------------------------
// GetTemplate cycle safety (fix 3.5)
// ---------------------------------------------------------------------------

func TestGetTemplateCyclicFather(t *testing.T) {
	a := &XTemplate{Name: "a", SubTemplates: make(map[string]*XTemplate)}
	b := &XTemplate{Name: "b", SubTemplates: make(map[string]*XTemplate)}
	// Create a cycle: a.Father = b, b.Father = a
	a.Father = b
	b.Father = a

	// GetTemplate must not infinite-loop; should return nil
	result := a.GetTemplate("nonexistent")
	if result != nil {
		t.Errorf("cyclic Father: expected nil, got non-nil")
	}
}

func TestGetTemplateFatherChain(t *testing.T) {
	root, _ := NewXTemplateFromString("&&nested&&")
	sub1, _ := NewXTemplateFromString("sub1-content")
	root.AddTemplate("nested", sub1)

	// sub1 can find "nested" via its Father (root)
	found := sub1.GetTemplate("nested")
	if found == nil {
		t.Error("GetTemplate via Father chain returned nil")
	}
}

// ---------------------------------------------------------------------------
// File size limit (fix 3.6)
// ---------------------------------------------------------------------------

func TestTemplateFileSizeLimit(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "tmpl_big_*.tmpl")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	defer tmp.Close()

	chunk := make([]byte, 1024*1024) // 1 MB
	for i := range chunk {
		chunk[i] = 'a'
	}
	for i := 0; i < 11; i++ { // 11 MB > 10 MB limit
		tmp.Write(chunk)
	}
	tmp.Close()

	_, err = NewXTemplateFromFile(tmp.Name())
	if err == nil {
		t.Error("oversized file should return error")
	}
}

// ---------------------------------------------------------------------------
// defer Close / ioutil removed (fix 3.8)
// ---------------------------------------------------------------------------

func TestLoadFileDeferClose(t *testing.T) {
	// Verifies that LoadFile works normally; defer close is implicitly tested
	// since a descriptor leak would only surface under extended load.
	// This test validates that a normal load+execute round-trip works correctly.
	tmpl, err := NewXTemplateFromFile("testunit/a.template")
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	data := XDataset{
		"clientname": "Fred",
		"hobbies": &XDatasetCollection{
			&XDataset{"name": "Football"},
			&XDataset{"name": "Ping-pong"},
		},
	}
	got := tmpl.Execute(&data)
	if got == "" {
		t.Error("LoadFile round-trip: empty output")
	}
}

// ---------------------------------------------------------------------------
// Compile errors
// ---------------------------------------------------------------------------

func TestCompileMismatchStart(t *testing.T) {
	_, err := NewXTemplateFromString("template [[subtemplate]]")
	if err == nil {
		t.Error("mismatch start/end should return error")
	}
}

func TestCompileMismatchEnd(t *testing.T) {
	_, err := NewXTemplateFromString("template [[]]")
	if err == nil {
		t.Error("extra [[]] should return error")
	}
}

func TestCompileNestedMismatch(t *testing.T) {
	_, err := NewXTemplateFromString("[[a]][[b]][[]]")
	if err == nil {
		t.Error("nested mismatch should return error")
	}
}

// ---------------------------------------------------------------------------
// Missing file
// ---------------------------------------------------------------------------

func TestNewXTemplateFromFileMissing(t *testing.T) {
	_, err := NewXTemplateFromFile("/nonexistent/path/template.tmpl")
	if err == nil {
		t.Error("missing file should return error")
	}
}

// ---------------------------------------------------------------------------
// Path separator in template names (pipe syntax)
// ---------------------------------------------------------------------------

func TestPipeSyntaxMultipleTemplates(t *testing.T) {
	// [[a|b]] registers the same content under both "a" and "b"
	tmpl, err := NewXTemplateFromString(`&&a&& &&b&&
[[a|b]]SHARED[[]]`)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	// The \n between &&b&& and [[a|b]] is outer-template text rendered after both references.
	got := tmpl.Execute(&XDataset{})
	if got != "SHARED SHARED\n" {
		t.Errorf("pipe template names: got %q, want 'SHARED SHARED\\n'", got)
	}
}

// ---------------------------------------------------------------------------
// String / GoString
// ---------------------------------------------------------------------------

func TestXTemplateStringEmpty(t *testing.T) {
	tmpl := &XTemplate{}
	if got := tmpl.String(); got != "xcore.XTemplate{}" {
		t.Errorf("String() on nil Root: got %q", got)
	}
	if got := tmpl.GoString(); got != "#xcore.XTemplate{}" {
		t.Errorf("GoString() on nil Root: got %q", got)
	}
}

func TestXTemplateStringNonEmpty(t *testing.T) {
	tmpl, _ := NewXTemplateFromString("hello")
	s := tmpl.String()
	if s == "" || s == "xcore.XTemplate{}" {
		t.Errorf("String() on compiled template: got %q", s)
	}
}

// ---------------------------------------------------------------------------
// Full complex template (b.template)
// ---------------------------------------------------------------------------

func TestXTemplateSimple(t *testing.T) {
	tmpl, err := NewXTemplateFromFile("testunit/b.template")
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	tmp, _ := time.Parse(time.RFC3339, "2020-01-01T12:00:00")
	lang := NewXLanguage("mainpage", language.English)
	lang.Set("welcome", "Welcome to you")
	data := XDataset{
		"clientname":    "Fred",
		"clientpicture": "face.jpg",
		"hobbies": &XDatasetCollection{
			&XDataset{"name": "Football", "sport": "yes"},
			&XDataset{"name": "Ping-pong", "sport": "yes"},
			&XDataset{"name": "Swimming", "sport": "yes"},
			&XDataset{"name": "Videogames", "sport": "no"},
			&XDataset{"name": "other 1", "sport": "no"},
			&XDataset{"name": "other 2", "sport": "no"},
			&XDataset{"name": "other 3", "sport": "yes"},
			&XDataset{"name": "other 4", "sport": "no"},
			&XDataset{"name": "other 5", "sport": nil},
		},
		"preferredhobby": &XDataset{"name": "Baseball", "sport": "yes"},
		"metadata": &XDataset{
			"preferred-color": "blue",
			"Salary":          3568.65,
			"hiredate":        tmp,
		},
		"#": lang,
	}
	if got := tmpl.Execute(&data); got == "" {
		t.Error("b.template: empty output")
	}
}

// ---------------------------------------------------------------------------
// Concurrent execution (thread safety once compiled — fix 3.10)
// ---------------------------------------------------------------------------

func TestConcurrentExecute(t *testing.T) {
	tmpl, _ := NewXTemplateFromString(`Hello {{name}}! @@items:item@@
[[item]]{{val}} [[]]
`)
	const goroutines = 50
	const iterations = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				data := XDataset{
					"name": fmt.Sprintf("user%d", id),
					"items": &XDatasetCollection{
						&XDataset{"val": fmt.Sprintf("item%d-%d", id, j)},
					},
				}
				result := tmpl.Execute(&data)
				if result == "" {
					t.Errorf("concurrent Execute: empty result (goroutine %d)", id)
				}
			}
		}(i)
	}
	wg.Wait()
}

func TestConcurrentExecuteSafe(t *testing.T) {
	tmpl, _ := NewXTemplateFromString("{{name}}")
	const goroutines = 40
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			data := XDataset{"name": fmt.Sprintf("n%d", id)}
			result, err := tmpl.ExecuteSafe(&data)
			if err != nil {
				t.Errorf("concurrent ExecuteSafe: unexpected error: %v", err)
			}
			if result != fmt.Sprintf("n%d", id) {
				t.Errorf("concurrent ExecuteSafe: got %q, want n%d", result, id)
			}
		}(i)
	}
	wg.Wait()
}

// ---------------------------------------------------------------------------
// File round-trip with output filters
// ---------------------------------------------------------------------------

func TestOutputFiltersInFile(t *testing.T) {
	// Write a template with output filters to a temp file and load it
	dir := t.TempDir()
	path := filepath.Join(dir, "filter.tmpl")
	content := `<p>{{title|html}}</p><script>var q="{{q|js}}";</script><a href="?q={{q|url}}">link</a>`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing template file: %v", err)
	}
	tmpl, err := NewXTemplateFromFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	data := XDataset{
		"title": `<script>alert("xss")</script>`,
		"q":     `hello world&more`,
	}
	got := tmpl.Execute(&data)
	// html filter: raw <script> tag must not appear in output
	if strings.Contains(got, "<script>alert") {
		t.Errorf("html filter not applied: %q", got)
	}
	// js filter: & must be encoded as \x26, NOT left as a bare &
	if strings.Contains(got, `world&more`) {
		t.Errorf("js filter not applied (& not escaped): %q", got)
	}
	if !strings.Contains(got, `world\x26more`) {
		t.Errorf("js filter not applied (\\x26 missing): %q", got)
	}
	// url filter: space must be encoded as + in the query string
	if !strings.Contains(got, "hello+world") {
		t.Errorf("url filter not applied: %q", got)
	}
}

// ---------------------------------------------------------------------------
// Existing regression tests (kept from original file)
// ---------------------------------------------------------------------------

func TestNewXTemplateFromString(t *testing.T) {
	tmpl, _ := NewXTemplateFromString(`
%-- This is a comment. It will not appear in the final code. --%
Let's put your name here: {{clientname}}<br />
And lets put your hobbies here:<br />
@@hobbies:hobby@@
[[hobby]]
I love {{name}}<br />
[[]]
`)
	data := XDataset{
		"clientname": "Fred",
		"hobbies": &XDatasetCollection{
			&XDataset{"name": "Football"},
			&XDataset{"name": "Ping-pong"},
			&XDataset{"name": "Swimming"},
			&XDataset{"name": "Videogames"},
		},
	}
	want := "\nLet's put your name here: Fred<br />\nAnd lets put your hobbies here:<br />\nI love Football<br />\nI love Ping-pong<br />\nI love Swimming<br />\nI love Videogames<br />\n\n"
	if got := tmpl.Execute(&data); got != want {
		t.Errorf("TestNewXTemplateFromString:\ngot:  %q\nwant: %q", got, want)
	}
}

func TestNewXTemplateFromFile(t *testing.T) {
	tmpl, _ := NewXTemplateFromFile("testunit/a.template")
	data := XDataset{
		"clientname": "Fred",
		"hobbies": &XDatasetCollection{
			&XDataset{"name": "Football"},
			&XDataset{"name": "Ping-pong"},
			&XDataset{"name": "Swimming"},
			&XDataset{"name": "Videogames"},
		},
	}
	want := "\nLet's put your name here: Fred<br />\nAnd lets put your hobbies here:<br />\nI love Football<br />\nI love Ping-pong<br />\nI love Swimming<br />\nI love Videogames<br />\n\n"
	if got := tmpl.Execute(&data); got != want {
		t.Errorf("TestNewXTemplateFromFile:\ngot:  %q\nwant: %q", got, want)
	}
}
