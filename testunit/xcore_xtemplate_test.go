package main

import (
	"fmt"
	"github.com/webability-go/xcore"
	"testing"
	//  "unsafe"
)

/* TEST XTEMPLATE */

func TestCommentParam(t *testing.T) {
	// Test 1: assign a simple parameter string with some comments
	tmpl := xcore.NewXTemplate()
	err := tmpl.LoadString("%-- starting comment --%Text %--with a [[]]comment--% here. Also an%----% empty comment %--ending comment--%")

	fmt.Println(err)
	//  fmt.Println(tmpl.Root)

	result := tmpl.Execute(nil)

	fmt.Println("Result: ", result)

	if result != "Text  here. Also an empty comment " {
		t.Errorf("The comments have not been removed correctly")
	}
}

func TestLanguageParam(t *testing.T) {
	tmpl, _ := xcore.NewXTemplateFromString("Test with ##some## ##languages## here")

	fmt.Println(tmpl)
	fmt.Println(tmpl.Root)

	data := &xcore.XDataset{}
	l, _ := xcore.NewXLanguageFromString("some=a tiny table\nlanguages=of english language\n")
	data.Set("#", l)

	result := tmpl.Execute(data)

	fmt.Println("Result: ", result)

	if result != "Test with a tiny table of english language here" {
		t.Errorf("The language table has not been inserted correctly")
	}
}
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
