package xcore

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/text/language"
)

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
	// The creation of the data is obviously tedious here, in real life it should come from a JSON, a Database, etc
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

func TestNewXTemplateFromString(t *testing.T) {
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
	// The creation of the data is obviously tedious here, in real life it should come from a JSON, a Database, etc
	data := XDataset{
		"clientname": "Fred",
		"hobbies": &XDatasetCollection{
			&XDataset{"name": "Football"},
			&XDataset{"name": "Ping-pong"},
			&XDataset{"name": "Swimming"},
			&XDataset{"name": "Videogames"},
		},
	}

	str := tmpl.Execute(&data)
	if str != `
Let's put your name here: Fred<br />
And lets put your hobbies here:<br />
I love Football<br />
I love Ping-pong<br />
I love Swimming<br />
I love Videogames<br />

` {
		t.Error("Error loading and running the template from string " + str)
		return
	}
}

func TestNewXTemplateFromFile(t *testing.T) {
	tmpl, _ := NewXTemplateFromFile("testunit/a.template")

	// The creation of the data is obviously tedious here, in real life it should come from a JSON, a Database, etc
	data := XDataset{
		"clientname": "Fred",
		"hobbies": &XDatasetCollection{
			&XDataset{"name": "Football"},
			&XDataset{"name": "Ping-pong"},
			&XDataset{"name": "Swimming"},
			&XDataset{"name": "Videogames"},
		},
	}

	str := tmpl.Execute(&data)
	if str != `
Let's put your name here: Fred<br />
And lets put your hobbies here:<br />
I love Football<br />
I love Ping-pong<br />
I love Swimming<br />
I love Videogames<br />

` {
		t.Error("Error loading and running the template from file " + str)
		return
	}
}

func TestXTemplateErrors(t *testing.T) {
	tmpl1, _ := NewXTemplateFromString("template [[subtemplate]]")

	if tmpl1 != nil {
		t.Error("Error compiling a template with error, the template should be nil")
		return
	}
}

func TestXTemplateComments(t *testing.T) {
	tmpl1, _ := NewXTemplateFromString("abcdefg")
	tmpl2, _ := NewXTemplateFromString(`a%--comment1--%b%--comment 2
[[]]
[[subtemplate]]
		--%c%--comment3 --%
defg%-- ending @@subtemplate@@ comment --%
`)

	if tmpl1.Execute(nil) != tmpl2.Execute(nil) {
		t.Error("Error comparing templates for comments")
		return
	}
}

func TestXTemplateLanguageParam(t *testing.T) {
	tmpl, _ := NewXTemplateFromString("Test with ##some## ##languages## here")

	data := &XDataset{}
	l, _ := NewXLanguageFromString("some=a tiny table\nlanguages=of english language\n")
	data.Set("#", l)

	result := tmpl.Execute(data)

	if result != "Test with a tiny table of english language here" {
		t.Errorf("The language table has not been inserted correctly")
	}
}

func TestXTemplateSimple(t *testing.T) {
	tmpl, err := NewXTemplateFromFile("testunit/b.template")
	if err != nil {
		t.Error(err)
		return
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
		},
		"preferredhobby": &XDataset{
			"name":  "Baseball",
			"sport": "yes",
		},
		"metadata": &XDataset{
			"preferred-color": "blue",
			"Salary":          3568.65,
			"hiredate":        tmp,
		},
		"#": lang,
	}

	str := tmpl.Execute(&data)
	if str == "" {
		t.Errorf("Error build complex template")
	}
}

func TestXTemplateClone(t *testing.T) {
	tmpl, err := NewXTemplateFromFile("testunit/b.template")
	if err != nil {
		t.Error(err)
		return
	}

	newtmpl := tmpl.Clone()

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
		},
		"preferredhobby": &XDataset{
			"name":  "Baseball",
			"sport": "yes",
		},
		"metadata": &XDataset{
			"preferred-color": "blue",
			"Salary":          3568.65,
			"hiredate":        tmp,
		},
		"#": lang,
	}

	str1 := tmpl.Execute(&data)
	if str1 == "" {
		t.Errorf("Error build complex template")
	}
	str2 := newtmpl.Execute(&data)
	if str2 == "" {
		t.Errorf("Error build complex template cloned")
	}
	if str2 != str1 {
		t.Errorf("Error comparing template cloned")
	}
	fmt.Println(str1)
}

/*
package main

import (
	"fmt"
	"github.com/webability-go/xcore"
	"testing"
	//  "unsafe"
)

// TEST XTEMPLATE

func TestReferenceParam(t *testing.T) {
	tmpl, _ := xcore.NewXTemplateFromString(`
The sub template starts here: &&template1&&. End.
[[template1]]
This is the template 1
[[]]
`)

	fmt.Println(tmpl)
	fmt.Println(tmpl.Root)

	result := tmpl.Execute(&xcore.XDataset{})

	fmt.Println("Result: ", result)

}

func TestComplexReferenceParam(t *testing.T) {

	tmpl, err := xcore.NewXTemplateFromString(`
The sub template starts here: &&template2&&. End.
[[template1]]
This is the template 1
[[]]
[[template2]]
This is the template 2
  [[template3]]
  This is the subtemplate 3
  [[]]
  [[template4|template5]]
  These are the subtemplates 4 and 5
    [[template6.first]]
    This is the subtemplate 6 first element for a loop
    [[]]
    [[template6]]
    This is the subtemplate 6 any element for a loop
    [[]]
    [[template6.last]]
    This is the subtemplate 6 last element for a loop
    [[]]
  [[]]
[[]]
[[template7|template7.status.false]]
This is the template 7 for field status false and any other values
[[]]
[[template7.status.true]]
This is the template 7 for field status true
[[]]
`)

	if err != nil {
		fmt.Println(err)
		return
	}

	result := tmpl.Execute(&xcore.XDataset{})
	fmt.Println("Result: ", result)
}

func TestVariableParam(t *testing.T) {

	tmpl, err := xcore.NewXTemplateFromString(`
Some data:
{{data1}}
{{data2}}
{{data3>data31}}
{{data4}}
{{data5}}
{{data6}}
{{data7}}
{{data8}}
@@data8@@
[[data8]]
* test {{data81}} and {{data82}} and {{data83}} and {{data1}}
[[]]
??data9??
[[data9]]
* Data 9 exists and is {{data9}}
[[]]
??data10??
[[data10]]
* Data 10 does not exist
[[]]
!!dump!!
`)

	if err != nil {
		fmt.Println(err)
		return
	}

	data := xcore.XDataset{}
	data["data1"] = "DATA1"
	data["data2"] = "DATA1"
	sm := xcore.XDataset{}
	sm["data31"] = "DATA31"
	data["data3"] = sm
	data["data4"] = 123
	data["data5"] = 123.432
	data["data6"] = true
	data["data7"] = func() string { return "ABC" }

	d8r1 := &xcore.XDataset{}
	d8r1.Set("data81", "rec 1: Entry 8-1")
	d8r1.Set("data82", "rec 1: Entry 8-2")

	d8r2 := &xcore.XDataset{}
	d8r2.Set("data81", "rec 2: Entry 8-1")
	d8r2.Set("data82", "rec 2: Entry 8-2")
	d8r2.Set("data83", "rec 2: Entry 8-3")

	d8r3 := &xcore.XDataset{}
	d8r3.Set("data81", "rec 3: Entry 8-1")
	d8r3.Set("data82", "rec 3: Entry 8-2")

	d := xcore.XDatasetCollection{}
	d.Push(d8r1)
	d.Push(d8r2)
	d.Push(d8r3)

	data["data8"] = &d
	data["data9"] = "I exist"

	fmt.Printf("Data: %v\n", data)
	//  fmt.Printf("ADDRESS DATA8 / GET R1: %p", data.GetCollection("data8").Get(0))

	result := tmpl.Execute(&data)
	fmt.Println("Result: ", result)
}
*/
