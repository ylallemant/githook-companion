package commit

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/nlp"
	nlpapi "github.com/ylallemant/githook-companion/pkg/nlp/api"
)

func Validate(message string, cfg *api.Config) (bool, string, *nlpapi.Dictionary) {
	validationRegexp := validationExpression(cfg)
	formatted := isMessageValid(message, validationRegexp)

	if !formatted {
		tokenizer, _ := nlp.NewTokenizer(cfg.Commit.TokenizerOptions)
		tokens, _, _ := tokenizer.Tokenize(message)

		fmt.Printf("Tokens :\n")
		for _, token := range tokens {
			fmt.Printf("   - token : %#+v\n", token)

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
