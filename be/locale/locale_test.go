package locale

import (
	"testing"

	"github.com/template/be/lib/language"
)

func TestLocale(t *testing.T) {
	testString := Translate(language.EN, "Test")
	testString2 := Translate(language.ID, "Test")
	testString3 := Translate("", "Test")

	if testString != "Test" {
		t.Error("Failed to translate English")
		return
	}
	if testString2 != "Tes" {
		t.Error("Failed to translate Bahasa")
		return
	}

	if testString3 != "" {
		t.Error("Failed to error handling locale")
		return
	}

}

func TestError(t *testing.T) {
	testString := Error(language.EN, "Test")
	testString2 := Error(language.ID, "Test")
	testString3 := Error("", "Test")

	if testString.Error() != "Test" {
		t.Error("Failed to translate English")
		return
	}
	if testString2.Error() != "Tes" {
		t.Error("Failed to translate Bahasa")
		return
	}

	if testString3.Error() != "" {
		t.Error("Failed to error handling locale")
		return
	}
}
