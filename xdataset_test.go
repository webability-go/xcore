package xcore

import (
	"fmt"
	"testing"
	"time"
)

func ExampleXDataset() {

	tmp, _ := time.Parse(time.RFC3339, "2020-01-01T12:00:00")
	ds := &XDataset{
		"v1":  123,
		"v2":  "abc",
		"vt":  tmp,
		"v3":  true,
		"vpi": 3.1415927,
	}
	fmt.Println(ds)

	data := &XDataset{
		"clientname":    "Fred",
		"clientpicture": "face.jpg",
		"hobbies": &XDatasetCollection{
			&XDataset{"name": "Football", "sport": "yes"},
			&XDataset{"name": "Ping-pong", "sport": "yes"},
			&XDataset{"name": "Swimming", "sport": "yes"},
			&XDataset{"name": "Videogames", "sport": "no"},
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
	}

	fmt.Println(data)
	// Output:
	// xcore.XDataset{v1:123 v2:abc v3:true vpi:3.1415927 vt:0001-01-01 00:00:00 +0000 UTC}
	// xcore.XDataset{clientname:Fred clientpicture:face.jpg hobbies:XDatasetCollection[0:xcore.XDataset{name:Football sport:yes} 1:xcore.XDataset{name:Ping-pong sport:yes} 2:xcore.XDataset{name:Swimming sport:yes} 3:xcore.XDataset{name:Videogames sport:no} ] metadata:xcore.XDataset{Salary:3568.65 hiredate:0001-01-01 00:00:00 +0000 UTC preferred-color:blue} preferredhobby:xcore.XDataset{name:Baseball sport:yes}}
}

func TestCreateSimpleXDataset(t *testing.T) {

	// 1. Create a simple XDataset
	tmp, _ := time.Parse(time.RFC3339, "2020-01-01T12:00:00.0Z")
	ds := &XDataset{
		"v1":  123,
		"v2":  "abc",
		"vt":  tmp,
		"v3":  true,
		"vpi": 3.1415927,
	}

	// 2. Set, Get and Del on simple
	fmt.Printf("STRING %v\n", ds)
	fmt.Printf("GOSTRING %#v\n", ds)
}

func TestCreateComplexXDataset(t *testing.T) {

	// 1. Create a complex XDataset with subsets

	tmp, _ := time.Parse(time.RFC3339, "2020-01-01T12:00:00")
	data := XDataset{
		"clientname":    "Fred",
		"clientpicture": "face.jpg",
		"hobbies": &XDatasetCollection{
			&XDataset{"name": "Football", "sport": "yes"},
			&XDataset{"name": "Ping-pong", "sport": "yes"},
			&XDataset{"name": "Swimming", "sport": "yes"},
			&XDataset{"name": "Videogames", "sport": "no"},
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
	}

	// 2. Set, Get and Del

	v1, _ := data.GetString("clientname")
	t.Log("clientname:", v1)
	v2, _ := data.Get("hobbies")
	t.Log("hobbies:", v2)
	v3, _ := data.Get("hobbies>2")
	t.Log("second set of hobbies:", v3)
	v4, _ := data.GetString("hobbies>2>name")
	t.Log("name of third set of hobbies:", v4)
}

func TestLoadJSONInXDataset(t *testing.T) {

}

func TestCloneXDataset(t *testing.T) {
	// please log
	LOG = true

	// Test 1: creates 100 max, no file, expires in 1 sec
	ds := &XDataset{
		"v1": 123,
		"v2": "abc",
		"v3": true,
		"v4": &XDataset{
			"v4_p1": 456,
			"v4_p2": "def",
			"v4_p3": false,
		},
	}
	fmt.Printf("ORIGINAL %#v\n", ds)

	cds := ds.Clone()
	(*(*ds)["v4"].(*XDataset))["p5"] = "val5"

	fmt.Println("CLONED", cds)

	dsc := &XDatasetCollection{
		&XDataset{"v1": 123,
			"v2": "abc",
			"v3": true,
			"v4": &XDataset{
				"v4_p1": 456,
				"v4_p2": "def",
				"v4_p3": false,
			},
		},
		&XDataset{"v11": 123,
			"v12": "abc",
			"v13": true,
			"v14": &XDataset{
				"v14_p1": 456,
				"v14_p2": "def",
				"v14_p3": false,
			},
		},
	}

	fmt.Printf("ORIGINAL %#v\n", *dsc)

	cdsc := dsc.Clone().(*XDatasetCollection)

	(*((*(*cdsc)[0].(*XDataset))["v4"].(*XDataset)))["algomas"] = "el valor mas"

	fmt.Println("CLONED", *cdsc)
	fmt.Printf("%v %v\n", ((*(*dsc)[0].(*XDataset))["v4"].(*XDataset)), ((*(*cdsc)[0].(*XDataset))["v4"].(*XDataset)))
}
