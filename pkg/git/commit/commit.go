package commit

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ylallemant/githook-companion/pkg/api"
)

const (
	commitTypePrefixRegexpFmt = "^(?i)%s\\s*:{0,1}\\s*"
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
