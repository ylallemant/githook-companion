package nlp

import (
	"math/big"
	"strings"

	"github.com/pemistahl/lingua-go"
	"github.com/pkg/errors"
)

const unknown = "unknown"

func NewLanguageDetector(languageCodes []string, threshold float64) (*LanguageDetector, error) {
	detector := new(LanguageDetector)
	detector.threshold = threshold

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
	threshold     float64
}

func (i *LanguageDetector) DetectLanguage(sentence string) (string, string, bool) {
	confidenceValues := i.detector.ComputeLanguageConfidenceValues(sentence)

	highestConfidence := confidenceValues[0].Value()

	bigFloatHighest := big.NewFloat(highestConfidence)
	bigFloatThreshold := big.NewFloat(i.threshold)
	result := bigFloatThreshold.Cmp(bigFloatHighest)

	if result > 0 {
		return unknown, unknown, false
	}

	languageCode := confidenceValues[0].Language().IsoCode639_1().String()
	languageName := confidenceValues[0].Language().String()

	return strings.ToLower(languageCode), strings.ToLower(languageName), true
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
