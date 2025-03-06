package nlp

import (
	"strings"

	"github.com/pemistahl/lingua-go"
	"github.com/pkg/errors"
)

const unknown = "unknown"

func NewLanguageDetector(languageCodes []string) (*LanguageDetector, error) {
	detector := new(LanguageDetector)
	languages := make([]lingua.Language, 0)

	for _, languageCode := range languageCodes {
		language, err := languageFromCode(languageCode)
		if err != nil {
			return nil, errors.Wrap(err, "failed to initiate language detector")
		}

		languages = append(languages, language)
	}

	detector.detector = lingua.NewLanguageDetectorBuilder().
		FromLanguages(languages...).
		WithMinimumRelativeDistance(0.1).
		Build()

	return detector, nil
}

type LanguageDetector struct {
	detector      lingua.LanguageDetector
	languageCodes []string
}

func (i *LanguageDetector) DetectLanguage(sentence string) (string, string, bool) {
	confidenceValues := i.detector.ComputeLanguageConfidenceValues(sentence)

	highestConfidence := confidenceValues[0]

	if highestConfidence.Value() < 0.8 {
		return unknown, unknown, false
	}

	code := strings.ToLower(highestConfidence.Language().IsoCode639_1().String())

	return code, strings.ToLower(highestConfidence.Language().String()), true
}

func languageFromCode(code string) (lingua.Language, error) {
	switch code {
	case "en":
		return lingua.English, nil
	case "de":
		return lingua.German, nil
	case "fr":
		return lingua.French, nil
	default:
		return lingua.Unknown, errors.Errorf("unknown language with code \"%s\"", code)
	}
}
