package commit

import (
	"regexp"
	"strings"
)

var (
	asciiRegexp      = regexp.MustCompile(`[^a-z0-9\s]`)
	whitespaceRegexp = regexp.MustCompile(`\s+`)
	acronymRegexp    = regexp.MustCompile(`(?P<first><!\w)([a-z])\.`)
)

func sanitiseString(text string) string {
	text = strings.ToLower(text)
	text = removeSpecialChars(text)
	return strings.TrimSpace(text)
}

func removeSpecialChars(text string) string {
	text = acronymRegexp.ReplaceAllString(text, `$1`)
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
