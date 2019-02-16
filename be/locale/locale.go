package locale

import (
	"errors"
	"fmt"

	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/template/be/lib/env"
	"github.com/template/be/lib/language"
)

type localeModule struct {
	TranslateID i18n.TranslateFunc
	TranslateEN i18n.TranslateFunc
}

var localem *localeModule

/*
Init locale module
Call this function first before using database module
*/
func Init(environment string) {

	if environment == env.Development {

		wdList := []string{"../../../files/var", "../../files/var", "../files/var", "./files/var"}
		for w := range wdList {
			if err := i18n.LoadTranslationFile(fmt.Sprintf("%s/locale/%s.all.json", wdList[w], language.ID)); err == nil {
				i18n.LoadTranslationFile(fmt.Sprintf("%s/locale/%s.all.json", wdList[w], language.EN))
				break
			}
		}

	} else {

		i18n.LoadTranslationFile(fmt.Sprintf("/var/locale/%s.all.json", language.ID))
		i18n.LoadTranslationFile(fmt.Sprintf("/var/locale/%s.all.json", language.EN))
	}

	indonesian, _ := i18n.Tfunc(language.ID)
	english, _ := i18n.Tfunc(language.EN)

	localem = &localeModule{
		TranslateID: indonesian,
		TranslateEN: english,
	}
}

// TFunc get translate function
func TFunc(lang string) i18n.TranslateFunc {
	if lang == language.EN {
		return localem.TranslateEN
	}

	return localem.TranslateID
}

// Translate using key
func Translate(lang, key string) string {

	if (lang != language.ID && lang != language.EN) || localem == nil {
		return ""
	}

	f := TFunc(lang)

	return f(key)

}

// Error result from translation
func Error(lang, key string) error {

	return errors.New(Translate(lang, key))

}
