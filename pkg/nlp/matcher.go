package nlp

import (
	"github.com/ylallemant/githook-companion/pkg/nlp/api"
)

func Match(text string, matcher *api.Matcher) bool {
	return matcher.Regex.MatchString(text)
}

func FindAll(text string, matcher *api.Matcher) []string {
	return matcher.Regex.FindAllString(text, -1)
}
