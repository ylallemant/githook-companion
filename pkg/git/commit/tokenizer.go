package commit

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	asciiRegexp      = regexp.MustCompile(`[^a-z0-9\s]`)
	whitespaceRegexp = regexp.MustCompile(`\s+`)
	acronymRegexp    = regexp.MustCompile(`\b(?:[a-zA-Z]\.){2,}([a-zA-Z]){0,1}`)
)

func sanitiseString(text string) string {
	text = strings.ToLower(text)
	text = removeDiacritics(text)
	text = replaceAcronyms(text)
	text = removeSpecialChars(text)
	return strings.TrimSpace(text)
}

func removeDiacritics(text string) string {
	transformer := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(transformer, text)
	return result
}

func replaceAcronyms(text string) string {
	acronyms := acronymRegexp.FindAll([]byte(text), -1)

	for _, acronym := range acronyms {
		cleaned := strings.ReplaceAll(string(acronym), ".", "")
		text = strings.ReplaceAll(text, string(acronym), cleaned)
	}

	return text
}

func removeSpecialChars(text string) string {
	text = asciiRegexp.ReplaceAllString(text, " ")
	text = whitespaceRegexp.ReplaceAllString(text, " ")

	return text
}

func spilt(text string) []string {
	return whitespaceRegexp.Split(text, -1)
}

func tonenize(text string) []string {
	text = sanitiseString(text)
	return spilt(text)
}
