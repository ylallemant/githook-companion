package commit

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ylallemant/githooks-butler/pkg/api"
	"github.com/ylallemant/githooks-butler/pkg/config"
)

func Validate(message string, cfg *api.Config) (bool, string, *api.CommitTypeDictionary) {
	validationRegexp := validationExpression(cfg)
	formatted := isMessageValid(message, validationRegexp)

	if !formatted {
		tokens := tonenize(message)

		dictionary := fuzzyDictionaryMatch(tokens[0], cfg)

		if dictionary != nil {
			return true, dictionary.Type, dictionary
		}

		return false, "", nil
	}

	commitType, _ := messageType(message, cfg)

	return true, commitType, nil
}

func validationExpression(cfg *api.Config) *regexp.Regexp {
	expression := fmt.Sprintf(
		commitTypePrefixRegexpFmt,
		strings.Join(config.GetCommitTypes(cfg), "|"),
	)

	formatted, _ := regexp.Compile(expression)

	return formatted
}

func isMessageValid(message string, validationRegexp *regexp.Regexp) bool {
	return validationRegexp.MatchString(message)
}
