package nlp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectLanguage(t *testing.T) {
	cases := []struct {
		name          string
		sentence      string
		expectedCode  string
		expectedName  string
		expectedKnown bool
	}{
		{
			name:          "known language german",
			sentence:      "neues Ladenkorb angelegt",
			expectedCode:  "de",
			expectedName:  "german",
			expectedKnown: true,
		},
		{
			name:          "known language english",
			sentence:      "fix wrong function output",
			expectedCode:  "en",
			expectedName:  "english",
			expectedKnown: true,
		},
		{
			name:          "unknown Klingon sentence",
			sentence:      "Heghlu'meH QaQ jajvam",
			expectedCode:  "unknown",
			expectedName:  "unknown",
			expectedKnown: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			detector, _ := NewLanguageDetector(DefaultLanguageDetectionOptions())

			code, name, known := detector.DetectLanguage(c.sentence, true)

			assert.Equal(tt, c.expectedCode, code, "wrong language code")
			assert.Equal(tt, c.expectedName, name, "wrong language name")
			assert.Equal(tt, c.expectedKnown, known, "wrong confidence")
		})
	}
}
