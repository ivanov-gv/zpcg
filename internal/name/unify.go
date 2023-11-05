package name

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Set func(r rune) bool

func (s Set) Contains(r rune) bool {
	return s(r)
}

var isMnOrSpace Set = func(r rune) bool {
	return unicode.Is(unicode.Mn, r) || // Mn: nonspacing marks
		r == ' '
}

func Unify(name string) string {
	t := transform.Chain(norm.NFD, runes.Remove(isMnOrSpace), norm.NFC)
	withoutUnicode, _, _ := transform.String(t, name)
	return strings.ToLower(withoutUnicode)
}
