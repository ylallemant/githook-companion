package nlp

import (
	"github.com/ylallemant/githook-companion/pkg/nlp/api"
)

const DefaultConfidenceThresthold = 0.6

func DefaultTokenizerOptions() *api.TokenizerOptions {
	options := new(api.TokenizerOptions)

	options.ConfidenceThresthold = DefaultConfidenceThresthold
	options.LanguageDetectionOptions = DefaultLanguageDetectionOptions()
	options.LanguageCodes = []string{
		"en",
	}

	return options
}
