package commit

import (
	"math"
	"slices"
	"strings"

	"github.com/ylallemant/githook-companion/pkg/api"
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
		if token.Source == nlpapi.TokenSourceLexeme && token.SourceName == api.CommitTypeTokenName {
			return token, true
		}
	}

	return nil, false
}

func assessMessageType(tokens []*nlpapi.Token, cfg *api.Config) (*nlpapi.Token, bool) {
	typeWeights := make(map[string]int)
	tokenMap := make(map[string]*nlpapi.Token)

	for _, token := range tokens {
		if token.Name != nlpapi.TokenUnknown && token.Source != nlpapi.TokenSourceLexeme {
			if _, ok := typeWeights[token.Name]; ok {
				typeWeights[token.SourceName] = typeWeights[token.SourceName] + commitTypeWeightIncrement(token)
			} else {
				typeWeights[token.SourceName] = commitTypeWeightIncrement(token)
				tokenMap[token.SourceName] = token
			}
		}
	}

	hierarchy := commitTypeHierarchy(cfg.Types)

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

	if !found {
		// no valid commit type was found in the tokens
		return nil, false
	}

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

func commitTypeWeightIncrement(token *nlpapi.Token) int {
	increment := 2

	if strings.Contains(token.SourceName, "weak") {
		increment = 1
	}

	if token.Source == nlpapi.TokenSourceLexeme && token.SourceName == api.CommitTypeTokenName {
		increment = 100000
	}

	return int(math.Floor(float64(increment) * token.Confidence))
}

func formatMessage(token *nlpapi.Token, template string) string {
	return ""
}
