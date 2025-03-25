package commit

import (
	"slices"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/nlp"
	nlpapi "github.com/ylallemant/githook-companion/pkg/nlp/api"
)

func Validate(message string, configuration *api.Config) (string, bool, *nlpapi.Token, []*nlpapi.Token, error) {
	tokenizer, _ := nlp.NewTokenizer(configuration.Commit.TokenizerOptions)
	languageCode, _, known := tokenizer.LanguageDetector().DetectLanguage(message, false)
	log.Debug().Msgf("detected language \"%s\"", languageCode)

	if !known {
		return languageCode, false, nil, []*nlpapi.Token{}, errors.Errorf("unknown language")
	}

	if !slices.Contains(configuration.Commit.TokenizerOptions.LanguageCodes, languageCode) {
		return languageCode, false, nil, []*nlpapi.Token{}, errors.Errorf("language detected in the commit message is not allowed (\"%s\")", languageCode)
	}

	tokens, languageCode, _, err := tokenizer.Tokenize(message)
	if err != nil {
		return languageCode, false, nil, tokens, errors.Wrap(err, "failed to tokenize commit message")
	}

	//validationRegexp := validationExpression(cfg)
	commitTypeToken, formatted := hasCommitTypeToken(tokens)
	log.Debug().Msgf("commit type token found: %v", commitTypeToken != nil)

	if !formatted {
		token, found := assessMessageType(tokens, configuration)

		if found {
			tokens = append(tokens, token)
			commitTypeToken = token
		}
	}

	if commitTypeToken == nil {
		return languageCode, false, nil, tokens, nil
	}

	return languageCode, true, commitTypeToken, tokens, nil
}
