package xcore

import (
	"fmt"
	"testing"
	"time"
)

func ExampleXDatasetCollectionTS() {

	tmp, _ := time.Parse(time.RFC3339, "2020-01-01T12:00:00.0Z")
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
	// xcore.XDataset{v1:123 v2:abc v3:true vpi:3.1415927 vt:2020-01-01 12:00:00 +0000 UTC}
	// xcore.XDataset{clientname:Fred clientpicture:face.jpg hobbies:XDatasetCollection[0:xcore.XDataset{name:Football sport:yes} 1:xcore.XDataset{name:Ping-pong sport:yes} 2:xcore.XDataset{name:Swimming sport:yes} 3:xcore.XDataset{name:Videogames sport:no} ] metadata:xcore.XDataset{Salary:3568.65 hiredate:2020-01-01 12:00:00 +0000 UTC preferred-color:blue} preferredhobby:xcore.XDataset{name:Baseball sport:yes}}
}

func TestXDatasetCollectionTS_simple_print(t *testing.T) {

	// 1. Create a simple XDataset
	tmp, _ := time.Parse(time.RFC3339, "2020-01-01T12:00:00.0Z")
	ds := &XDataset{
		"v1":  123,
		"v2":  "abc",
		"vt":  tmp,
		"v3":  true,
		"vpi": 3.1415927,
	}

	// 2. print
	str := fmt.Sprintf("%v", ds)
	if str != "xcore.XDataset{v1:123 v2:abc v3:true vpi:3.1415927 vt:2020-01-01 12:00:00 +0000 UTC}" {
		t.Error("Error creating and printing simple XDataset " + str)
		return
	}

	str = fmt.Sprintf("%#v", ds)
	if str != "#xcore.XDataset{v1:123 v2:\"abc\" v3:true vpi:3.1415927 vt:time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC)}" {
		t.Error("Error creating and #printing simple XDataset " + str)
		return
	}
}
