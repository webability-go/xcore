package xcore

import (
	"fmt"
	"testing"
	"time"
)

func ExampleXDataset() {

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

func TestXDataset_simple_print(t *testing.T) {

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
	if str != "#xcore.XDataset{v1:123 v2:\"abc\" v3:true vpi:3.1415927 vt:time.Time{wall:0x0, ext:63713476800, loc:(*time.Location)(nil)}}" {
		t.Error("Error creating and #printing simple XDataset " + str)
		return
	}
}

func getComplexDataset() *XDataset {

	// Create a complex XDataset with subsets
	tmp, _ := time.Parse(time.RFC3339, "2020-01-01T12:00:00")
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
	return data
}

func TestCreateXDataset_complex_print(t *testing.T) {

	data := getComplexDataset()

	str := fmt.Sprintf("%v", data)
	if str != "xcore.XDataset{clientname:Fred clientpicture:face.jpg hobbies:XDatasetCollection[0:xcore.XDataset{name:Football sport:yes} 1:xcore.XDataset{name:Ping-pong sport:yes} 2:xcore.XDataset{name:Swimming sport:yes} 3:xcore.XDataset{name:Videogames sport:no} ] metadata:xcore.XDataset{Salary:3568.65 hiredate:0001-01-01 00:00:00 +0000 UTC preferred-color:blue} preferredhobby:xcore.XDataset{name:Baseball sport:yes}}" {
		t.Error("Error creating and printing complex XDataset " + str)
		return
	}

	str = fmt.Sprintf("%#v", data)
	if str != "#xcore.XDataset{clientname:\"Fred\" clientpicture:\"face.jpg\" hobbies:XDatasetCollection[0:xcore.XDataset{name:Football sport:yes} 1:xcore.XDataset{name:Ping-pong sport:yes} 2:xcore.XDataset{name:Swimming sport:yes} 3:xcore.XDataset{name:Videogames sport:no} ] metadata:#xcore.XDataset{Salary:3568.65 hiredate:time.Time{wall:0x0, ext:0, loc:(*time.Location)(nil)} preferred-color:\"blue\"} preferredhobby:#xcore.XDataset{name:\"Baseball\" sport:\"yes\"}}" {
		t.Error("Error creating and #printing complex XDataset " + str)
		return
	}
}

func TestCreateXDataset_Get(t *testing.T) {

	data := getComplexDataset()

	// 2. Gets and paths
	v1, _ := data.GetString("clientname")
	if v1 != "Fred" {
		t.Error("Error getting clientname from XDataset " + v1)
		return
	}

	v2, _ := data.Get("hobbies")
	strv2 := fmt.Sprintf("%v", v2)
	if strv2 != "XDatasetCollection[0:xcore.XDataset{name:Football sport:yes} 1:xcore.XDataset{name:Ping-pong sport:yes} 2:xcore.XDataset{name:Swimming sport:yes} 3:xcore.XDataset{name:Videogames sport:no} ]" {
		t.Error("Error getting hobbies from XDataset " + strv2)
		return
	}

	v3, _ := data.Get("hobbies>2")
	strv3 := fmt.Sprintf("%v", v3)
	if strv3 != "xcore.XDataset{name:Swimming sport:yes}" {
		t.Error("Error getting hobbies>2 from XDataset " + strv3)
		return
	}

	v4, _ := data.GetString("hobbies>2>name")
	if v4 != "Swimming" {
		t.Error("Error getting hobbies>2>name from XDataset " + v4)
		return
	}
}

func TestLoadJSONInXDataset(t *testing.T) {

}

func TestXDataset_Clone(t *testing.T) {

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
	original := fmt.Sprintf("%v", ds)

	cds := ds.Clone()
	cloned := fmt.Sprintf("%v", cds)

	if original != cloned {
		t.Error("Error during clonation of XDataset " + cloned)
		return
	}

	// Verify it has been really cloned
	(*(*ds)["v4"].(*XDataset))["p5"] = "val5"
	val5o, _ := ds.GetString("v4>p5")  // should be "val5"
	val5c, _ := cds.GetString("v4>p5") // should be ""
	if val5o == val5c {
		t.Error("Error during clonation of XDataset " + val5c)
		return
	}
}

func TestXDatasetCollection_Clone(t *testing.T) {

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

	original := fmt.Sprintf("%v", dsc)

	cdsc := dsc.Clone().(*XDatasetCollection)
	cloned := fmt.Sprintf("%v", cdsc)

	if original != cloned {
		t.Error("Error during clonation of XDatasetCollection " + cloned)
		return
	}

	// Verify it has been really cloned
	(*((*(*dsc)[0].(*XDataset))["v4"].(*XDataset)))["p5"] = "val5"

	ro0, _ := dsc.Get(0)
	val5o, _ := ro0.GetString("v4>p5") // should be "val5"

	r0, _ := cdsc.Get(0)
	val5c, _ := r0.GetString("v4>p5") // should be ""

	if val5c != "" || val5o != "val5" {
		t.Error("Error during clonation of XDatasetCollection " + val5c)
		return
	}
}
