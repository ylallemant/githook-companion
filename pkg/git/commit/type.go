package commit

import (
	"slices"

	"github.com/rs/zerolog/log"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/nlp"
	nlpapi "github.com/ylallemant/githook-companion/pkg/nlp/api"
)

func commitTypeHierarchy(types []*api.CommitType) []string {
	hierarchy := make([]string, 0)

	for _, current := range types {
		hierarchy = append(hierarchy, current.Type)
	}

	return hierarchy
}

func hasCommitTypeToken(tokens []*nlpapi.Token) (*nlpapi.Token, bool) {
	for _, token := range tokens {
		if (token.Source == nlpapi.TokenSourceLexeme || token.Source == nlpapi.TokenSourceLexemeComposite) && token.SourceName == api.CommitTypeTokenName {
			return token, true
		}
	}

	return nil, false
}

func assessMessageType(tokens []*nlpapi.Token, cfg *api.Config) (*nlpapi.Token, bool) {
	log.Debug().Msgf("assess commit type from %d tokens", len(tokens))
	typeWeights := make(map[string]int)
	tokenMap := make(map[string]*nlpapi.Token)

	for _, token := range tokens {
		if token.Name != nlpapi.TokenUnknown && token.Source != nlpapi.TokenSourceLexeme && token.Source != nlpapi.TokenSourceLexemeComposite {
			dictionary := nlp.DictionaryByName(token.SourceName, cfg.TokenizerOptions)
			if dictionary == nil {
				log.Warn().Msgf("dictionary from token source name \"%s\" was not found", token.SourceName)
				continue
			}

			if _, ok := typeWeights[token.Name]; ok {
				typeWeights[token.SourceName] = typeWeights[token.SourceName] + dictionary.Weight
			} else {
				typeWeights[token.SourceName] = dictionary.Weight
				tokenMap[token.SourceName] = token
			}
		}
	}
	log.Debug().Msgf("commit type weights: %v", typeWeights)

	hierarchy := commitTypeHierarchy(cfg.Types)
	log.Debug().Msgf("commit type hierarchy: %v", hierarchy)

	var highestWeightToken *nlpapi.Token
	highestWeight := 0
	found := false

	for tokenSourceName, weight := range typeWeights {
		if highestWeight < weight {
			highestWeightToken = tokenMap[tokenSourceName]
			highestWeight = weight
			found = true
		} else if found && highestWeight == weight {
			if commitTypeHierarchyNumber(hierarchy, highestWeightToken.Value) > commitTypeHierarchyNumber(hierarchy, tokenMap[tokenSourceName].Value) {
				highestWeightToken = tokenMap[tokenSourceName]
				highestWeight = weight
			}
		}
	}
	log.Debug().Msgf("commit type token found: %v", found)

	if !found {
		// no valid commit type was found in the tokens
		return nil, false
	}

	log.Debug().Msgf("result: token \"%s\" with value \"%s\"", highestWeightToken.Name, highestWeightToken.Value)

	return commitTypeTokenFromToken(highestWeightToken), found
}

func commitTypeTokenFromToken(token *nlpapi.Token) *nlpapi.Token {
	return &nlpapi.Token{
		Name:        api.CommitTypeTokenName,
		Value:       token.Value,
		Source:      token.Source,
		SourceName:  token.SourceName,
		SourceMatch: token.SourceMatch,
		Word:        token.Word,
		Confidence:  token.Confidence,
		Index:       -1,
	}
}

func CommitTypeTokenFromString(commitType, languageCode string) *nlpapi.Token {
	return &nlpapi.Token{
		Name:        api.CommitTypeTokenName,
		Value:       commitType,
		Source:      nlpapi.TokenSourceLexeme,
		SourceName:  "forced-value",
		SourceMatch: commitType,
		Word: &nlpapi.Word{
			LanguageCode: nlpapi.LanguageCodeWildcard,
			Raw:          commitType,
			Cleaned:      commitType,
			Normalised:   commitType,
			Source:       nlpapi.TokenSourceLexeme,
			SourceName:   "forced-value",
		},
		Confidence: 1,
		Index:      -1,
	}
}

func commitTypeHierarchyNumber(hierarchy []string, typeName string) int {
	number := slices.Index(hierarchy, typeName)

	if number < 0 {
		number = 10000
	}

	return number
}
