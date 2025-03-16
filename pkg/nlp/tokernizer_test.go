package nlp

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ylallemant/githook-companion/pkg/nlp/api"
)

func TestNewTokenizer(t *testing.T) {
	cases := []struct {
		name                 string
		options              *api.TokenizerOptions
		word                 *api.Word
		expectError          bool
		expectedErrorMessage string
	}{
		{
			name:        "empty options",
			options:     &api.TokenizerOptions{},
			expectError: false,
		},
		{
			name: "single language",
			options: &api.TokenizerOptions{
				LanguageCodes: []string{"en"},
			},
			expectError: false,
		},
		{
			name: "empty dictionary",
			options: &api.TokenizerOptions{
				Dictionaries: []*api.Dictionary{},
			},
			expectError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			tokenizer, err := NewTokenizer(c.options)

			if c.expectError {
				assert.Nil(tt, tokenizer)
				assert.NotNil(tt, err)
				if err != nil {
					assert.Equal(tt, c.expectedErrorMessage, err.Error(), "wrong error massage")
				}
			} else {
				assert.NotNil(tt, tokenizer)
				assert.Nil(tt, err)
			}
		})
	}
}

func TestTokenizer_fuzzyDictionaryMatch(t *testing.T) {
	cases := []struct {
		name                   string
		options                *api.TokenizerOptions
		word                   *api.Word
		expectMatch            bool
		expectedDictionaryName string
		expectedEntry          string
		expectedConfidence     float64
	}{
		{
			name:    "no dictionary, no match",
			options: &api.TokenizerOptions{},
			word: &api.Word{
				LanguageCode: "de",
				Raw:          "Döner",
				Cleaned:      "doner",
				Normalised:   "doner",
			},
			expectedDictionaryName: "",
			expectedEntry:          "",
		},
		{
			name: "ignore empty dictionary",
			options: &api.TokenizerOptions{
				Dictionaries: []*api.Dictionary{},
			},
			word: &api.Word{
				LanguageCode: "de",
				Raw:          "Döner",
				Cleaned:      "doner",
				Normalised:   "doner",
			},
			expectedDictionaryName: "",
			expectedEntry:          "",
		},
		{
			name: "ignore foreign dictionary",
			options: &api.TokenizerOptions{
				Dictionaries: []*api.Dictionary{
					{
						LanguageCode: "en",
						Name:         "dishes",
						TokenName:    "dish",
						Entries: []string{
							"doner",
						},
					},
				},
			},
			word: &api.Word{
				LanguageCode: "de",
				Raw:          "Döner",
				Cleaned:      "doner",
				Normalised:   "doner",
			},
			expectedDictionaryName: "",
			expectedEntry:          "",
		},
		{
			name: "simple match",
			options: &api.TokenizerOptions{
				Dictionaries: []*api.Dictionary{
					{
						LanguageCode: "de",
						Name:         "Gerichte",
						TokenName:    "dish",
						Entries: []string{
							"donor",
							"doner",
						},
					},
				},
			},
			word: &api.Word{
				LanguageCode: "de",
				Raw:          "Döner",
				Cleaned:      "doner",
				Normalised:   "doner",
			},
			expectMatch:            true,
			expectedDictionaryName: "Gerichte",
			expectedEntry:          "doner",
			expectedConfidence:     1,
		},
		{
			name: "partial match",
			options: &api.TokenizerOptions{
				Dictionaries: []*api.Dictionary{
					{
						LanguageCode: "de",
						Name:         "Dönerwelt",
						TokenName:    "dish",
						Entries: []string{
							"donerbude",
							"donermann",
						},
					},
				},
			},
			word: &api.Word{
				LanguageCode: "de",
				Raw:          "Döner",
				Cleaned:      "doner",
				Normalised:   "doner",
			},
			expectMatch:            true,
			expectedDictionaryName: "Dönerwelt",
			expectedEntry:          "donerbude",
			expectedConfidence:     0.7777777777777778,
		},
		{
			name: "below global confidence threshold",
			options: &api.TokenizerOptions{
				ConfidenceThresthold: 0.9,
				Dictionaries: []*api.Dictionary{
					{
						LanguageCode: "de",
						Name:         "Dönerwelt",
						TokenName:    "dish",
						Entries: []string{
							"donerbude",
							"donermann",
						},
					},
				},
			},
			word: &api.Word{
				LanguageCode: "de",
				Raw:          "Döner",
				Cleaned:      "doner",
				Normalised:   "doner",
			},
			expectMatch:            false,
			expectedDictionaryName: "",
			expectedEntry:          "",
			expectedConfidence:     0,
		},
		{
			name: "below dictionary confidence threshold",
			options: &api.TokenizerOptions{
				ConfidenceThresthold: 0.6,
				Dictionaries: []*api.Dictionary{
					{
						ConfidenceThresthold: 0.9,
						LanguageCode:         "de",
						Name:                 "Dönerwelt",
						TokenName:            "dish",
						Entries: []string{
							"donerbude",
							"donermann",
						},
					},
				},
			},
			word: &api.Word{
				LanguageCode: "de",
				Raw:          "Döner",
				Cleaned:      "doner",
				Normalised:   "doner",
			},
			expectMatch:            false,
			expectedDictionaryName: "",
			expectedEntry:          "",
			expectedConfidence:     0,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			tokenizer, err := NewTokenizer(c.options)
			assert.Nil(tt, err)

			dictionary, entry, confidence := tokenizer.fuzzyDictionaryMatch(c.word)

			if c.expectMatch {
				assert.NotNil(tt, dictionary)
				assert.Equal(tt, c.expectedEntry, entry, "wrong entry")
				assert.Equal(tt, c.expectedConfidence, confidence, "wrong confidence")
				if dictionary != nil {
					assert.Equal(tt, c.expectedDictionaryName, dictionary.Name, "wrong dictionary name")
				}
			} else {
				assert.Nil(tt, dictionary)
				assert.Equal(tt, c.expectedEntry, entry, "wrong entry")
				assert.Equal(tt, c.expectedConfidence, confidence, "wrong confidence")
			}
		})
	}
}

func TestTokenizer_Tokenize(t *testing.T) {
	cases := []struct {
		name                 string
		options              *api.TokenizerOptions
		sentence             string
		expectedLanguage     string
		expectedTemplate     string
		expectedTokens       []*api.Token
		expectError          bool
		expectedErrorMessage string
	}{
		{
			name:             "empty sentence",
			sentence:         "",
			options:          &api.TokenizerOptions{},
			expectedTemplate: "",
			expectedTokens:   []*api.Token{},
		},
		{
			name:     "ignore empty dictionary",
			sentence: "(GDT-3564) added new, cool, feature",
			options: &api.TokenizerOptions{
				ConfidenceThresthold: 0.8,
				Lexemes: []*api.Lexeme{
					{
						LanguageCode: api.LanguageCodeWildcard,
						TokenName:    "issue-tracker-reference",
						Variants: []*api.Variant{
							{
								Matcher: &api.Matcher{Regex: regexp.MustCompile("[\\(\\[]{0,1}([\\w]{0,6})[-_]([\\d]+)[\\)\\]]{0,1}")},
								Normalisers: []*api.NormalisationStep{
									{
										Matcher:    &api.Matcher{Regex: regexp.MustCompile("([\\w]{0,6})[-_]([\\d]+)")},
										ReplaceAll: true,
										Formatter: &api.Formatter{
											Template: "{{ upper . }}",
										},
									},
									{
										Matcher:     &api.Matcher{Regex: regexp.MustCompile("[-_]")},
										Replacement: "-",
									},
								},
							},
							{
								Matcher: &api.Matcher{Regex: regexp.MustCompile("[\\(\\[]{0,1}(#|gh-|GH-)([\\d]+)[\\)\\]]{0,1}")},
								Normalisers: []*api.NormalisationStep{
									{
										Name:        "github issue reference",
										Matcher:     &api.Matcher{Regex: regexp.MustCompile("(#)")},
										Replacement: "gh-",
									},
								},
							},
						},
						Normalisers: []*api.NormalisationStep{
							{
								Matcher:    &api.Matcher{Regex: regexp.MustCompile("(#|gh-|GH-)([\\d]+)")},
								ReplaceAll: true,
								Formatter: &api.Formatter{
									Template: "{{ lower . }}",
								},
							},
						},
					},
				},
				Dictionaries: []*api.Dictionary{
					{
						LanguageCode:      "en",
						Name:              "Git commit type feature",
						TokenName:         "feat",
						TokenValueIsMatch: true,
						Entries: []string{
							"add",
							"implement",
							"feature",
							"feat",
						},
					},
				},
			},
			expectedLanguage: "en",
			expectedTemplate: "word~0 word~1 word~2, word~3, word~4",
			expectedTokens: []*api.Token{
				{
					Name:       "issue-tracker-reference",
					Source:     api.TokenSourceLexeme,
					Value:      "GDT-3564",
					Confidence: 1,
				},
				{
					Name:       "feat",
					Source:     api.TokenSourceDictionary,
					Value:      "add",
					Confidence: 1,
				},
				{
					Name:       api.TokenUnknown,
					Source:     api.TokenSourceNone,
					Value:      "new",
					Confidence: 0,
				},
				{
					Name:       api.TokenUnknown,
					Source:     api.TokenSourceNone,
					Value:      "cool",
					Confidence: 0,
				},
				{
					Name:       "feat",
					Source:     api.TokenSourceDictionary,
					Value:      "feature",
					Confidence: 1,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			tokenizer, err := NewTokenizer(c.options)
			assert.Nil(tt, err)

			tokens, languageCode, template, err := tokenizer.Tokenize(c.sentence)

			if c.expectError {
				assert.Equal(tt, c.expectedTemplate, template, "wrong template")
				assert.Equal(tt, c.expectedTokens, tokens, "wrong tokens")

				assert.NotNil(tt, err)
				if err != nil {
					assert.Equal(tt, c.expectedErrorMessage, err.Error(), "wrong error massage")
				}
			} else {
				assert.Nil(tt, err)
				assert.Equal(tt, c.expectedTemplate, template, "wrong template")
				assert.Equal(tt, c.expectedLanguage, languageCode, "wrong language code")
				for idx, token := range tokens {
					assert.Equal(tt, c.expectedTokens[idx].Name, token.Name, "wrong tokens")
					assert.Equal(tt, c.expectedTokens[idx].Source, token.Source, "wrong tokens")
					assert.Equal(tt, c.expectedTokens[idx].Confidence, token.Confidence, "wrong tokens")
					assert.Equal(tt, c.expectedTokens[idx].Value, token.Value, "wrong tokens")
				}
			}
		})
	}
}
