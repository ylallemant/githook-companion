package commit

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ylallemant/githooks-butler/pkg/api"
)

const (
	commitTypePrefixRegexpFmt = "^(?i)%s\\s*:{0,1}\\s*"
)

var (
	fistWord = regexp.MustCompile(`^\w+`)
)

func messageType(message string, cfg *api.Config) (string, bool) {
	message = strings.ToLower(message)

	for _, commitType := range cfg.Commit.Types {
		expression, _ := regexp.Compile(fmt.Sprintf(
			commitTypePrefixRegexpFmt,
			commitType.Type,
		))

		if expression.MatchString(message) {
			return strings.ToLower(commitType.Type), true
		}
	}

	return "", false
}

func EnsureFormat(message, commitType string) string {
	expression, _ := regexp.Compile(fmt.Sprintf(
		commitTypePrefixRegexpFmt,
		commitType,
	))

	if expression.MatchString(message) {
		// format commit type prefix
		return expression.ReplaceAllString(
			message,
			fmt.Sprintf("%s: ", strings.ToLower(commitType)),
		)
	} else {
		// add commit type prefix to message
		return fmt.Sprintf(
			"%s: %s",
			strings.ToLower(commitType),
			message,
		)
	}
}

func EnsureDictionaryValue(message string, dictionary *api.CommitTypeDictionary) string {
	message = strings.TrimSpace(message)

	return fistWord.ReplaceAllString(message, dictionary.Value)
}
