package xcore

import (
	"fmt"
	"golang.org/x/text/language"
	"testing"
)

func ExampleXLanguage() {

	lang_es := NewXLanguage("messages", language.Spanish)
	lang_en := NewXLanguage("messages", language.English)
	lang_fr := NewXLanguage("messages", language.French)

	// You can load this from your system files for instance and keep your translation tables apart
	lang_es.Set("panicerror", "Error crítico del sistema")
	lang_en.Set("panicerror", "System panic error")
	lang_fr.Set("panicerror", "Erreur grave dans le système")

	// pick a random general system language (for instance the client's OS language)
	lang := lang_fr

	// launch a panic errors
	// panic(lang.Get("panicerror"))
	fmt.Println("Launch a panic error message in the selected language:", lang.Get("panicerror"))
}

func TestXLanguage(t *testing.T) {

	lang_es := NewXLanguage("messages", language.Spanish)
	lang_en := NewXLanguage("messages", language.English)
	lang_fr := NewXLanguage("messages", language.French)

	// You can load this from your system files for instance and keep your translation tables apart
	lang_es.Set("panicerror", "Error crítico del sistema")
	lang_en.Set("panicerror", "System panic error")
	lang_fr.Set("panicerror", "Erreur grave dans le système")

	// pick a random general system language (for instance the client's OS language)
	lang := lang_fr

	// launch a panic errors
	// panic(lang.Get("panicerror"))
	t.Log("Launch a panic error message in the selected language:", lang.Get("panicerror"))
}
