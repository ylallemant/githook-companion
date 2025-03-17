package commit

import (
	"fmt"

	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/nlp"
	nlpapi "github.com/ylallemant/githook-companion/pkg/nlp/api"
)

// TODO: return error too
func Validate(message string, cfg *api.Config) (string, bool, *nlpapi.Token, []*nlpapi.Token) {
	tokenizer, _ := nlp.NewTokenizer(cfg.Commit.TokenizerOptions)

	// TODO: get dictionary from function to test its structure
	tokenizer.AddDictionary(&nlpapi.Dictionary{
		LanguageCode:      nlpapi.LanguageCodeWildcard,
		Name:              fmt.Sprintf("%s_dictionary", api.CommitTypeTokenName),
		TokenName:         api.CommitTypeTokenName,
		TokenValueIsMatch: true,
		Entries:           config.CommitTypeList(cfg),
	})

	tokens, languageCode, _, _ := tokenizer.Tokenize(message)

	//validationRegexp := validationExpression(cfg)
	commitTypeToken, formatted := hasCommitTypeToken(tokens)

	if !formatted {
		token, found := assessMessageType(tokens, cfg)

		if found {
			tokens = append(tokens, token)
			commitTypeToken = token
		}
	}

	if commitTypeToken == nil {
		return languageCode, false, nil, tokens
	}

	return languageCode, true, commitTypeToken, tokens
}
