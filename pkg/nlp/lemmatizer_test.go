package nlp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ylallemant/githook-companion/pkg/nlp/api"
)

func TestLemmatizer(t *testing.T) {
	cases := []struct {
		name          string
		langCode      string
		expectError   bool
		expectedError string
	}{
		{
			name:        "known language",
			langCode:    "en",
			expectError: false,
		},
		{
			name:          "faile on unknown language",
			langCode:      "klingon",
			expectError:   true,
			expectedError: "unsupported language \"klingon\"",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			lemmatizer, err := Lemmatizer(c.langCode)

			if c.expectError {
				assert.Nil(tt, lemmatizer)
				assert.NotNil(tt, err)
				assert.Equal(tt, c.expectedError, err.Error(), "wrong error massage")
			} else {
				assert.NotNil(tt, lemmatizer)
				assert.Nil(tt, err)
			}
		})
	}
}

func TestLemmatizer_Lemma(t *testing.T) {
	cases := []struct {
		name          string
		langCode      string
		word          string
		expected      string
		expectError   bool
		expectedError string
	}{
		{
			name:     "resolved word lemmatization",
			langCode: "en",
			word:     "added",
			expected: "add",
		},
		{
			name:     "return input if word can't be lemmatized",
			langCode: "en",
			word:     "alacrite",
			expected: "alacrite",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			word := new(api.Word)
			word.LanguageCode = c.langCode
			word.Raw = c.word
			word.Cleaned = c.word

			lemmatizer, err := Lemmatizer(c.langCode)
			lemmatizer.Lemma(word)

			if c.expectError {
				assert.Nil(tt, lemmatizer)
				assert.NotNil(tt, err)
				assert.Equal(tt, c.expectedError, err.Error(), "wrong error massage")
			} else {
				assert.NotNil(tt, lemmatizer)
				assert.Nil(tt, err)
				assert.Equal(tt, c.expected, word.Normalised, "wrong error massage")
			}
		})
	}
}
