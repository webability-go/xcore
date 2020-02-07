package xcore

import (
	"fmt"
	"testing"
)

func ExampleNewXTemplateFromString() {
	tmpl, _ := NewXTemplateFromString(`
%-- This is a comment. It will not appear in the final code. --%
Let's put your name here: {{clientname}}<br />
And lets put your hobbies here:<br />
@@hobbies:hobby@@     %-- note the 1rst id is the entry into the data to inject and the second one is the name of the sub-template to use --%

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
}

func TestNewXTemplateFromString(t *testing.T) {
	tmpl, _ := NewXTemplateFromString(`
%-- This is a comment. It will not appear in the final code. --%
Let's put your name here: {{clientname}}<br />
And lets put your hobbies here:<br />
@@hobbies:hobby@@     %-- note the 1rst id is the entry into the data to inject and the second one is the name of the sub-template to use --%

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
}
