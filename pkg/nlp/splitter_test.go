package nlp

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ylallemant/githook-companion/pkg/nlp/api"
)

func TestSplitter_clean(t *testing.T) {
	cases := []struct {
		name             string
		sentence         string
		languageCode     string
		lexemes          []*api.Lexeme
		expectedSentence string
		expectedWords    []*api.Word
	}{
		{
			name:             "no lexeme no change",
			sentence:         "neues - schöneres - Döner Shop, implementiert !",
			languageCode:     "de",
			lexemes:          []*api.Lexeme{},
			expectedSentence: "neues  schöneres  Döner Shop implementiert",
			expectedWords:    []*api.Word{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			splitter := DefaultSplitter(c.languageCode, c.lexemes)

			sentence := splitter.clean(c.sentence)

			assert.Equal(tt, c.expectedSentence, sentence, "wrong sentence")
		})
	}
}

func TestSplitter_LanguageCode(t *testing.T) {
	cases := []struct {
		name         string
		languageCode string
		lexemes      []*api.Lexeme
	}{
		{
			name:         "no lexeme no change",
			languageCode: "de",
			lexemes:      []*api.Lexeme{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			splitter := DefaultSplitter(c.languageCode, c.lexemes)

			assert.Equal(tt, c.languageCode, splitter.LanguageCode(), "wrong sentence")
		})
	}
}

func TestSplitter_ExtractLexemes(t *testing.T) {
	cases := []struct {
		name             string
		sentence         string
		languageCode     string
		lexemes          []*api.Lexeme
		expectedSentence string
		expectedWords    map[string]*api.Word
	}{
		{
			name:             "no lexeme no change",
			sentence:         "neues - schöneres - Döner Shop implementiert !",
			languageCode:     "de",
			lexemes:          []*api.Lexeme{},
			expectedSentence: "neues - schöneres - Döner Shop implementiert !",
			expectedWords:    map[string]*api.Word{},
		},
		{
			name:         "wildcard lexeme",
			sentence:     "neues Döner Shop implementiert (gh-2345)",
			languageCode: "de",
			lexemes: []*api.Lexeme{
				{
					LanguageCode: api.LanguageCodeWildcard,
					TokenName:    "issue-tracker-reference",
					Variants: []*api.Variant{
						{
							Matcher: &api.Matcher{Regex: regexp.MustCompile("[\\(\\[]{0,1}([\\w]{0,6}[-_][\\d]+)[\\)\\]]{0,1}")},
						},
					},
				},
			},
			expectedSentence: "neues Döner Shop implementiert lexeme:0",
			expectedWords: map[string]*api.Word{
				"lexeme:0": {
					LanguageCode: api.LanguageCodeWildcard,
					FromLexeme:   "issue-tracker-reference",
					Raw:          "(gh-2345)",
					Cleaned:      "(gh-2345)",
					Normalised:   "(gh-2345)",
				},
			},
		},
		{
			name:         "wildcard lexeme",
			sentence:     "neues Döner Shop implementiert (#789) and [ECOM_2345]",
			languageCode: "de",
			lexemes: []*api.Lexeme{
				{
					LanguageCode: api.LanguageCodeWildcard,
					TokenName:    "issue-tracker-reference",
					Variants: []*api.Variant{
						{
							Matcher: &api.Matcher{Regex: regexp.MustCompile("[\\(\\[]{0,1}([\\w]{0,6})[-_]([\\d]+)[\\)\\]]{0,1}")},
						},
						{
							Matcher: &api.Matcher{Regex: regexp.MustCompile("[\\(\\[]{0,1}(#[\\d]+)[\\)\\]]{0,1}")},
						},
					},
					Normalisers: []*api.NormalisationStep{
						{
							Matcher:    &api.Matcher{Regex: regexp.MustCompile("(#|[\\w]{0,6})?([-_])?([\\d]+)")},
							ReplaceAll: true,
						},
						{
							Matcher:     &api.Matcher{Regex: regexp.MustCompile("[-_]")},
							Replacement: "-",
						},
					},
				},
			},
			expectedSentence: "neues Döner Shop implementiert lexeme:1 and lexeme:0",
			expectedWords: map[string]*api.Word{
				"lexeme:0": {
					LanguageCode: api.LanguageCodeWildcard,
					FromLexeme:   "issue-tracker-reference",
					Raw:          "[ECOM_2345]",
					Cleaned:      "ECOM-2345",
					Normalised:   "ECOM-2345",
				},
				"lexeme:1": {
					LanguageCode: api.LanguageCodeWildcard,
					FromLexeme:   "issue-tracker-reference",
					Raw:          "(#789)",
					Cleaned:      "#789",
					Normalised:   "#789",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			splitter := DefaultSplitter(c.languageCode, c.lexemes)

			sentence, words := splitter.ExtractLexemes(c.sentence)

			assert.Equal(tt, c.expectedSentence, sentence, "wrong sentence")
			assert.Equal(tt, c.expectedWords, words, "wrong words")
		})
	}
}

func TestSplitter_Split(t *testing.T) {
	cases := []struct {
		name             string
		sentence         string
		languageCode     string
		lexemes          []*api.Lexeme
		expectedTemplate string
		expectedWords    []*api.Word
	}{
		{
			name:             "unchanges words",
			sentence:         " neues Döner Shop, implementiert !\n",
			languageCode:     "de",
			lexemes:          []*api.Lexeme{},
			expectedTemplate: "word:0 word:1 word:2, word:3 !",
			expectedWords: []*api.Word{
				{
					LanguageCode: "de",
					Raw:          "neues",
				},
				{
					LanguageCode: "de",
					Raw:          "Döner",
				},
				{
					LanguageCode: "de",
					Raw:          "Shop",
				},
				{
					LanguageCode: "de",
					Raw:          "implementiert",
				},
			},
		},
		{
			name:         "wildcard lexeme",
			sentence:     "neues Döner Shop implementiert (#789) and [ecom_2345]",
			languageCode: "de",
			lexemes: []*api.Lexeme{
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
			expectedTemplate: "word:0 word:1 word:2 word:3 word:4 word:5 word:6",
			expectedWords: []*api.Word{
				{
					LanguageCode: "de",
					Raw:          "neues",
				},
				{
					LanguageCode: "de",
					Raw:          "Döner",
				},
				{
					LanguageCode: "de",
					Raw:          "Shop",
				},
				{
					LanguageCode: "de",
					Raw:          "implementiert",
				},
				{
					LanguageCode: api.LanguageCodeWildcard,
					FromLexeme:   "issue-tracker-reference",
					Raw:          "(#789)",
					Cleaned:      "gh-789",
					Normalised:   "gh-789",
				},
				{
					LanguageCode: "de",
					Raw:          "and",
				},
				{
					LanguageCode: api.LanguageCodeWildcard,
					FromLexeme:   "issue-tracker-reference",
					Raw:          "[ecom_2345]",
					Cleaned:      "ECOM-2345",
					Normalised:   "ECOM-2345",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			splitter := DefaultSplitter(c.languageCode, c.lexemes)

			template, words := splitter.Split(c.sentence)

			assert.Equal(tt, c.expectedTemplate, template, "wrong template")
			assert.Equal(tt, c.expectedWords, words, "wrong words")
		})
	}
}
