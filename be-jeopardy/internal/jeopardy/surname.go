package jeopardy

import (
	_ "embed"
	"strings"
	"unicode"
)

//go:embed first_names.txt
var firstNamesList string

var firstNames = toSet(strings.Fields(firstNamesList))

// Lowercase words allowed inside a name, as in "Vincent van Gogh" or "John of Gaunt".
var nameParticles = toSet(strings.Fields("van von de da di du la le del della der den ten bin al el y of"))

// Trailing parts dropped before taking the surname, so "Elizabeth I" doesn't accept "I".
var nameSuffixes = toSet(strings.Fields("jr jr. sr sr. i ii iii iv v vi vii viii"))

// A name-shaped answer ending in one of these is a place or thing, not a person:
// "Marshall Plan", "George Washington Bridge", "Charles River".
var nonPersonLastWords = toSet(strings.Fields(`
	act airplane airport avenue award awards band bay beach boulevard bowl bridge brothers
	building canal cape cathedral center church city code college company county crater
	cup dam day desert disease doctrine effect expedition falls forest foundation garden
	gardens gate glacier gulf hall harbor highway hole hotel hour house inc institute
	island islands lake law line memorial motors mountain mountains museum ocean orchestra
	park peak plan point prize quartet records rebellion river road rule scale school
	sea show sisters sound square stadium station street strait syndrome test theater
	theatre theorem tower trail treaty trio trophy university valley war way
`))

func toSet(words []string) map[string]bool {
	set := make(map[string]bool, len(words))
	for _, w := range words {
		set[strings.ToLower(w)] = true
	}
	return set
}

// surnameForms returns the shorter forms of a person's name that also count as
// correct: every run of words after the first name, e.g. "Vincent van Gogh" gives
// "van Gogh" and "Gogh". It returns nothing unless the answer looks like a
// person's name: two to five words, starting with a known first name, every other
// word capitalized or a name particle, and not ending in a place or thing word.
func surnameForms(answer string) []string {
	words := strings.Fields(answer)
	for len(words) > 0 && nameSuffixes[strings.ToLower(strings.TrimRight(words[len(words)-1], ","))] {
		words = words[:len(words)-1]
	}
	if len(words) < 2 || len(words) > 5 {
		return nil
	}
	if !firstNames[strings.ToLower(trimPunct(words[0]))] {
		return nil
	}
	for _, w := range words {
		w = trimPunct(w)
		if w == "" {
			return nil
		}
		if !unicode.IsUpper([]rune(w)[0]) && !nameParticles[strings.ToLower(w)] {
			return nil
		}
	}
	last := strings.ToLower(trimPunct(words[len(words)-1]))
	if nonPersonLastWords[last] || nameParticles[last] {
		return nil
	}

	forms := []string{}
	for i := 1; i < len(words); i++ {
		form := strings.TrimRight(strings.Join(words[i:], " "), ",")
		// "(William Randolph) Hearst" would give "Randolph) Hearst". The
		// alternatives already hold the form without the parentheses.
		if strings.ContainsAny(form, "()") || len([]rune(trimPunct(form))) < 3 {
			continue
		}
		forms = append(forms, form)
	}
	return forms
}

func trimPunct(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}
