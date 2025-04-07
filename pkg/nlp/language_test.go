package nlp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectLanguage(t *testing.T) {
	cases := []struct {
		name            string
		sentence        string
		ignoreWordCount bool
		expectedCode    string
		expectedName    string
		expectedKnown   bool
	}{
		{
			name:            "short german sentence",
			sentence:        "neue Ladenkorb",
			ignoreWordCount: true,
			expectedCode:    "de",
			expectedName:    "german",
			expectedKnown:   true,
		},
		{
			name:            "known language german",
			sentence:        "neues Ladenkorb angelegt",
			ignoreWordCount: true,
			expectedCode:    "de",
			expectedName:    "german",
			expectedKnown:   true,
		},
		{
			name:            "known language english",
			sentence:        "fix wrong function output",
			ignoreWordCount: true,
			expectedCode:    "en",
			expectedName:    "english",
			expectedKnown:   true,
		},
		{
			name:            "unknown Klingon sentence",
			sentence:        "Heghlu'meH QaQ jajvam",
			ignoreWordCount: true,
			expectedCode:    "unknown",
			expectedName:    "unknown",
			expectedKnown:   false,
		},
		{
			name:            "default language on bad word count",
			sentence:        "Heghlu'meH QaQ jajvam",
			ignoreWordCount: false,
			expectedCode:    "en",
			expectedName:    "english",
			expectedKnown:   true,
		},
		{
			name:            "long message",
			sentence:        "add dictionary weigth to make commit-type assertion configurable",
			ignoreWordCount: false,
			expectedCode:    "en",
			expectedName:    "english",
			expectedKnown:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			detector, _ := NewLanguageDetector(DefaultLanguageDetectionOptions())

			code, name, known := detector.DetectLanguage(c.sentence, c.ignoreWordCount)

			assert.Equal(tt, c.expectedCode, code, "wrong language code")
			assert.Equal(tt, c.expectedName, name, "wrong language name")
			assert.Equal(tt, c.expectedKnown, known, "wrong confidence")
		})
	}
}
