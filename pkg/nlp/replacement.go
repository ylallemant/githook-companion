package nlp

import (
	"fmt"
	"regexp"
)

func secureReplaceAllString(text, old, new string) string {
	regex := searchStringToRegex(old)

	return regex.ReplaceAllString(text, new)
}

// TODO: add error return type
func searchStringToRegex(search string) *regexp.Regexp {
	if search == "" {
		return nil
	}

	lastIndex := len(search) - 1
	escaped := fmt.Sprintf("(%s)", regexp.QuoteMeta(search))

	if !nonScriptCharRegexp.MatchString(search[0:1]) {
		escaped = fmt.Sprintf("\\b%s", escaped)
	}

	if !nonScriptCharRegexp.MatchString(search[lastIndex:]) {
		escaped = fmt.Sprintf("%s\\b", escaped)
	}

	regex, _ := regexp.Compile(escaped)

	return regex
}
