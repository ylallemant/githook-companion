package nlp

import (
	"strings"

	"github.com/pemistahl/lingua-go"
)

const unknown = "unknown"

func DetectLanguage(sentence string) (string, string, bool) {
	detector := lingua.NewLanguageDetectorBuilder().
		FromAllSpokenLanguages().
		WithMinimumRelativeDistance(0.2).
		Build()

	language, known := detector.DetectLanguageOf(sentence)

	code := strings.ToLower(language.IsoCode639_1().String())

	if !known {
		code = unknown
	}

	return code, strings.ToLower(language.String()), known
}
