package commit

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ylallemant/githook-companion/pkg/api"
	nlpapi "github.com/ylallemant/githook-companion/pkg/nlp/api"
)

func TestIsMessageValid(t *testing.T) {
	cases := []struct {
		name       string
		message    string
		commitType string
		expected   bool
	}{
		{
			name:       "no matching commit type prefix",
			message:    "some message",
			commitType: "feat",
			expected:   false,
		},
		{
			name:       "matching commit type prefix",
			message:    "feat: some message",
			commitType: "feat",
			expected:   true,
		},
		{
			name:       "ignore prefix case",
			message:    "FEAT: some message",
			commitType: "feat",
			expected:   true,
		},
		{
			name:       "ignore missing colon",
			message:    "FEAT some message",
			commitType: "feat",
			expected:   true,
		},
		{
			name:       "ignore missing spaces",
			message:    "FEAT:some message",
			commitType: "feat",
			expected:   true,
		},
		{
			name:       "ignore to mutch space",
			message:    "FEAT :   some message",
			commitType: "feat",
			expected:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			validationExpression, err := regexp.Compile(
				fmt.Sprintf(commitTypePrefixRegexpFmt, c.commitType),
			)
			assert.Nil(tt, err)

			result := isMessageValid(c.message, validationExpression)

			assert.Equal(tt, c.expected, result, "wrong result")
		})
	}
}

func TestValidationExpression(t *testing.T) {
	cases := []struct {
		name     string
		config   *api.Config
		expected string
	}{
		{
			name: "no matching commit type prefix",
			config: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{
						{
							Type: "feat",
						},
					},
				},
			},
			expected: "^(?i)feat\\s*:{0,1}\\s*",
		},
		{
			name: "no types declared",
			config: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{},
				},
			},
			expected: "^(?i)\\s*:{0,1}\\s*",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			result := validationExpression(c.config)

			expectedRegexp, err := regexp.Compile(c.expected)
			assert.Nil(tt, err)

			assert.Equal(tt, expectedRegexp, result, "regular expressions do not match")
		})
	}
}

func TestValidate(t *testing.T) {
	cases := []struct {
		name                     string
		message                  string
		config                   *api.Config
		expectedTokens           []*nlpapi.Token
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
								TokenName:    "feat",
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
								Name:         "add",
								TokenName:    "feat",
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
								TokenName:    "feat",
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
			valid, commitType, dictionary := Validate(c.message, c.config)

			if c.expectedDictionaryFound {
				assert.Equal(tt, c.config.Commit.TokenizerOptions.Dictionaries[0], dictionary, "wrong result")
			} else {
				assert.Nil(tt, dictionary)
			}

			assert.Equal(tt, c.expectedCommitType, commitType, "wrong result")
			assert.Equal(tt, c.expectedValidationResult, valid, "wrong result")
		})
	}
}
