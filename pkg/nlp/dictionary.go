package nlp

import (
	"github.com/ylallemant/githook-companion/pkg/nlp/api"
)

func DictionaryByName(name string, configuration *api.TokenizerOptions) *api.Dictionary {
	for _, dictionary := range configuration.Dictionaries {
		if dictionary.Name == name {
			return dictionary
		}
	}
	return nil
}
