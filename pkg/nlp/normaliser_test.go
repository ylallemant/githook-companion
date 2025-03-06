package nlp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ylallemant/githook-companion/pkg/nlp/api"
)

func Test_removeSpecialChars(t *testing.T) {
	cases := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "no lexeme no change",
			text:     "neues - schöneres - Döner Shop, implementiert !",
			expected: "neues sch neres ner hop implementiert ",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			result := removeSpecialChars(c.text)
			assert.Equal(tt, c.expected, result, "wrong result")
		})
	}
}

func Test_replaceAcronyms(t *testing.T) {
	cases := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "replace acronym",
			text:     "an der Uni kann man B.A.F.o.G. bekommen",
			expected: "an der Uni kann man BAFoG bekommen",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			result := replaceAcronyms(c.text)
			assert.Equal(tt, c.expected, result, "wrong result")
		})
	}
}

func Test_removeDiacritics(t *testing.T) {
	cases := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "no lexeme no change",
			text:     "neues - schöneres - Döner Shop, implementiert !",
			expected: "neues - schoneres - Doner Shop, implementiert !",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			result := removeDiacritics(c.text)
			assert.Equal(tt, c.expected, result, "wrong result")
		})
	}
}

func TestNormaliser_LanguageCode(t *testing.T) {
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
			normaliser, err := DefaultNormaliser(c.languageCode)
			assert.Nil(tt, err)

			assert.Equal(tt, c.languageCode, normaliser.LanguageCode(), "wrong sentence")
		})
	}
}

func TestNormaliser_Clean(t *testing.T) {
	cases := []struct {
		name         string
		languageCode string
		word         *api.Word
		expectedWord *api.Word
	}{
		{
			name:         "remove diacritics",
			languageCode: "en",
			word: &api.Word{
				LanguageCode: "en",
				Raw:          "Döner",
			},
			expectedWord: &api.Word{
				LanguageCode: "en",
				Raw:          "Döner",
				Cleaned:      "doner",
			},
		},
		{
			name:         "add lemma",
			languageCode: "en",
			word: &api.Word{
				LanguageCode: "en",
				Raw:          "adding",
			},
			expectedWord: &api.Word{
				LanguageCode: "en",
				Raw:          "adding",
				Cleaned:      "adding",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			normaliser, err := DefaultNormaliser(c.languageCode)
			assert.Nil(tt, err)

			normaliser.Clean(c.word)

			assert.Equal(tt, c.expectedWord, c.word, "wrong words")
		})
	}
}

func TestNormaliser_Normalise(t *testing.T) {
	cases := []struct {
		name         string
		languageCode string
		word         *api.Word
		expectedWord *api.Word
	}{
		{
			name:         "remove diacritics",
			languageCode: "en",
			word: &api.Word{
				LanguageCode: "en",
				Raw:          "Döner",
			},
			expectedWord: &api.Word{
				LanguageCode: "en",
				Raw:          "Döner",
				Cleaned:      "doner",
				Normalised:   "doner",
			},
		},
		{
			name:         "add lemma",
			languageCode: "en",
			word: &api.Word{
				LanguageCode: "en",
				Raw:          "added",
			},
			expectedWord: &api.Word{
				LanguageCode: "en",
				Raw:          "added",
				Cleaned:      "added",
				Normalised:   "add",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			normaliser, err := DefaultNormaliser(c.languageCode)
			assert.Nil(tt, err)

			normaliser.Normalise(c.word)

			assert.Equal(tt, c.expectedWord, c.word, "wrong words")
		})
	}
}

func TestNormaliser_NormaliseAll(t *testing.T) {
	cases := []struct {
		name          string
		languageCode  string
		words         []*api.Word
		expectedWords []*api.Word
	}{
		{
			name:         "unchanges words",
			languageCode: "de",
			words: []*api.Word{
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
					Raw:          "[ECOM_2345]",
					Cleaned:      "ECOM-2345",
					Normalised:   "ECOM-2345",
				},
			},
			expectedWords: []*api.Word{
				{
					LanguageCode: "de",
					Raw:          "neues",
					Cleaned:      "neues",
					Normalised:   "neu",
				},
				{
					LanguageCode: "de",
					Raw:          "Döner",
					Cleaned:      "doner",
					Normalised:   "doner",
				},
				{
					LanguageCode: "de",
					Raw:          "Shop",
					Cleaned:      "shop",
					Normalised:   "shop",
				},
				{
					LanguageCode: "de",
					Raw:          "implementiert",
					Cleaned:      "implementiert",
					Normalised:   "implementiert",
				},
				{
					LanguageCode: api.LanguageCodeWildcard,
					FromLexeme:   "issue-tracker-reference",
					Raw:          "[ECOM_2345]",
					Cleaned:      "ECOM-2345",
					Normalised:   "ECOM-2345",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			normaliser, err := DefaultNormaliser(c.languageCode)
			assert.Nil(tt, err)

			normaliser.NormaliseAll(c.words)

			assert.Equal(tt, c.expectedWords, c.words, "wrong words")
		})
	}
}
