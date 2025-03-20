package commit

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/nlp"
	nlpapi "github.com/ylallemant/githook-companion/pkg/nlp/api"
)

func TestHasCommitTypeToken(t *testing.T) {
	cases := []struct {
		name       string
		message    string
		commitType string
		expected   bool
	}{
		{
			name:     "no matching commit type prefix",
			message:  "some message",
			expected: false,
		},
		{
			name:       "matching commit type prefix",
			message:    "feat: some message",
			commitType: "FEAT",
			expected:   true,
		},
		{
			name:       "ignore prefix case",
			message:    "FEAT: some message",
			commitType: "FEAT",
			expected:   true,
		},
		{
			name:       "ignore missing colon",
			message:    "FEAT some message",
			commitType: "FEAT",
			expected:   true,
		},
		{
			name:       "ignore missing spaces",
			message:    "FEAT:some message",
			commitType: "FEAT",
			expected:   true,
		},
		{
			name:       "ignore to mutch space",
			message:    "FEAT :   some message",
			commitType: "FEAT",
			expected:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			defaultConfig := config.Default()
			tokenizer, err := nlp.NewTokenizer(defaultConfig.Commit.TokenizerOptions)
			assert.Nil(tt, err)

			tokens, _, _, err := tokenizer.Tokenize(c.message)
			assert.Nil(tt, err)

			token, result := hasCommitTypeToken(tokens)

			assert.Equal(tt, c.expected, result, "wrong result")

			if c.expected {
				assert.NotNil(tt, token)

				if token != nil {
					assert.Equal(tt, c.commitType, token.Value, "wrong result")
				}
			} else {
				assert.Nil(tt, token)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	cases := []struct {
		name                     string
		message                  string
		config                   *api.Config
		expectedToken            *nlpapi.Token
		expectedDictionaryFound  bool
		expectedCommitType       string
		expectedValidationResult bool
	}{
		{
			name:    "no matching commit type prefix",
			message: "some changes",
			config: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{
						{
							Type: "feat",
						},
					},
					TokenizerOptions: &nlpapi.TokenizerOptions{
						LanguageCodes: []string{
							"en",
						},
						Dictionaries: []*nlpapi.Dictionary{},
					},
				},
			},
			expectedDictionaryFound:  false,
			expectedCommitType:       "",
			expectedValidationResult: false,
		},
		{
			name:    "matching commit type prefix",
			message: "FeaT  :  some changes",
			config: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{
						{
							Type: "feat",
						},
					},
					TokenizerOptions: &nlpapi.TokenizerOptions{
						LanguageCodes: []string{
							"en",
						},
						Dictionaries: []*nlpapi.Dictionary{
							{
								LanguageCode: "en",
								Name:         "add",
								TokenName:    api.CommitTypeTokenName,
								TokenValue:   "feat",
								Entries: []string{
									"adds",
									"added",
									"adding",
									"new",
								},
							},
						},
					},
				},
			},
			expectedToken: &nlpapi.Token{
				Name:        api.CommitTypeTokenName,
				Source:      nlpapi.TokenSourceDictionary,
				SourceName:  fmt.Sprintf("%s_dictionary", api.CommitTypeTokenName),
				SourceMatch: "feat",
				Value:       "feat",
				Confidence:  1,
			},
			expectedDictionaryFound:  false,
			expectedCommitType:       "feat",
			expectedValidationResult: true,
		},
		{
			name:    "matching commit type with a dictionary",
			message: "added some changes",
			config: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{
						{
							Type: "feat",
						},
					},
					TokenizerOptions: &nlpapi.TokenizerOptions{
						LanguageCodes: []string{
							"en",
						},
						Dictionaries: []*nlpapi.Dictionary{
							{
								LanguageCode: "en",
								Name:         "weak-feature-signals",
								TokenName:    api.CommitTypeTokenName,
								TokenValue:   "feat",
								Entries: []string{
									"add",
									"implement",
									"use",
									"new",
								},
							},
						},
					},
				},
			},
			expectedToken: &nlpapi.Token{
				Name:        api.CommitTypeTokenName,
				Source:      nlpapi.TokenSourceDictionary,
				SourceName:  "weak-feature-signals",
				SourceMatch: "add",
				Value:       "feat",
				Confidence:  1,
			},
			expectedDictionaryFound:  true,
			expectedCommitType:       "feat",
			expectedValidationResult: true,
		},
		{
			name:    "no matching commit type with a dictionary",
			message: "added some changes",
			config: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{
						{
							Type: "feat",
						},
					},
					TokenizerOptions: &nlpapi.TokenizerOptions{
						ConfidenceThresthold: 0.8,
						LanguageCodes: []string{
							"en",
						},
						Dictionaries: []*nlpapi.Dictionary{
							{
								LanguageCode: "en",
								Name:         "nomatch",
								TokenName:    api.CommitTypeTokenName,
								TokenValue:   "feat",
								Entries: []string{
									"nomatch",
									"whatever",
								},
							},
						},
					},
				},
			},
			expectedDictionaryFound:  false,
			expectedCommitType:       "",
			expectedValidationResult: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			languageCode, valid, commitTypeToken, _, _ := Validate(c.message, c.config)

			if c.expectedToken != nil {
				assert.NotNil(tt, commitTypeToken)
				assert.Equal(tt, "en", languageCode, "wrong language code")

				if commitTypeToken != nil {
					assert.Equal(tt, c.expectedToken.Name, commitTypeToken.Name, "wrong token name")
					assert.Equal(tt, c.expectedToken.Source, commitTypeToken.Source, "wrong token source")
					assert.Equal(tt, c.expectedToken.SourceName, commitTypeToken.SourceName, "wrong token source name")
					assert.Equal(tt, c.expectedToken.SourceMatch, commitTypeToken.SourceMatch, "wrong token source match")
					assert.Equal(tt, c.expectedToken.Value, commitTypeToken.Value, "wrong token value")
					assert.Equal(tt, c.expectedToken.Confidence, commitTypeToken.Confidence, "wrong token confidence")
				}
			} else {
				assert.Nil(tt, commitTypeToken)
			}

			assert.Equal(tt, c.expectedValidationResult, valid, "wrong result")
		})
	}
}
