package xcore

import (
	"fmt"
	"golang.org/x/text/language"
	"testing"
	"time"
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

	t.Log(tmpl.Execute(&data))
}

func TestXTemplateComments(t *testing.T) {
	tmpl, _ := NewXTemplateFromString(`
%-- This is a comment. It will not appear in the final code. --%
Let's put your name here: {{clientname}}<br />
And lets put your hobbies here:<br />
@@hobbies:hobby@@     %-- note the 1rst id is the [[]] {{}} {{abc}} [[123]] entry into the data to inject and the second one is the name of the sub-template to use --%

%-- And you need the template for each hobby:--%
[[hobby]]I love {{name}}<br />
[[]]
%-- commented meta
[[hobby]]nothing[[]]
@@data@@
??exp??
--%
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

	t.Log(tmpl.Execute(&data))
}

func TestXTemplateSimple(t *testing.T) {
	tmpl, err := NewXTemplateFromString(`
&&header&&
&&body&&

[[header]]
Sports shop:
[[]]

[[body]]
{{clientname}} Preferred hobby:
&&:preferredhobby>sport:sport.&&  %-- will build sport_ + [yes/no] contained into the sport field. Be sure you have a template for each value ! --%

??preferredhobby>sport:sport??

[[sport.yes]]{{preferredhobby>name}} - It's a sport, sell him things![[]]
[[sport.no]]{{preferredhobby>name}} - It's not a sport, recommend him next store.[[]]
[[sport]]{{preferredhobby>name}} - We do not know that it is.[[]]



@@hobbies:hobby@@
[[hobby.first]]
1. {{name}} {{sport}}
[[]]
[[hobby.last]]
last. {{name}} {{sport}}
[[]]
[[hobby.even]]
2x. {{name}} {{sport}}
[[]]
[[hobby.key.3]]
3. {{name}} {{sport}}
[[]]
[[hobby]]
{{name}} {{sport}}
[[]]
[[hobby.none]]
No hobbies
[[]]
[[]]
  `)
	if err != nil {
		t.Errorf(err)
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

	t.Log(tmpl.Execute(&data))
}
