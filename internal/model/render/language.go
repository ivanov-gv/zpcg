package render

import (
	"github.com/samber/lo"
	"golang.org/x/text/language"
)

const BelarusianLanguageCode = "be"

var Belarusian = lo.Must(language.Parse(BelarusianLanguageCode)) // there is no var language.Belarusian, so we have to improvise

var DefaultLanguageTag = language.English

var SupportedLanguages = []language.Tag{
	language.Russian,
	language.Ukrainian,
	Belarusian,
	language.English,
	language.German,
	language.Serbian,
	language.Croatian,
	language.Slovak,
	language.Turkish,
}
