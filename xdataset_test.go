package xcore

import (
	"fmt"
	"testing"
	"time"
)

func ExampleXDataset() {

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

	fmt.Println(data)
}

func TestXDataset(t *testing.T) {

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

	v1, _ := data.GetString("clientname")
	fmt.Println("clientname:", v1)
	v2, _ := data.GetString("hobbies")
	fmt.Println("hobbies:", v2)
	v3, _ := data.GetString("hobbies>2")
	fmt.Println("second set of hobbies:", v3)
	v4, _ := data.GetString("hobbies>2>name")
	fmt.Println("name of third set of hobbies:", v4)
}
