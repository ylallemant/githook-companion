package commit

import (
	"fmt"
	"strings"

	"github.com/ylallemant/githooks-butler/pkg/api"
)

func Validate(message string, config *api.Config) (string, bool) {
	tokens := tonenize(message)

	dictionary := fuzzyDictionaryMatch(tokens[0], config)

	if dictionary != nil {
		tokens[0] = dictionary.Value
		message = fmt.Sprintf("%s: %s", dictionary.Type, strings.Join(tokens, " "))
		return message, true
	}

	return message, false
}
