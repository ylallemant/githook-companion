package nlp

import (
	"math/big"
	"strings"

	"github.com/pemistahl/lingua-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/ylallemant/githook-companion/pkg/nlp/api"
)

const unknown = "unknown"

func NewLanguageDetector(options *api.LanguageDetectionOptions) (api.LanguageDetector, error) {
	if options == nil {
		options = DefaultLanguageDetectionOptions()
	}

	detector := new(languageDetector)
	detector.threshold = options.ConfidenceThresthold
	detector.defaultLanguageCode = options.DefautLanguageCode
	detector.defaultLanguageName = options.DefautLanguageName

	detector.minimumWordCount = options.MinimumWordCount

	detector.detector = lingua.NewLanguageDetectorBuilder().
		FromAllLanguages().
		WithPreloadedLanguageModels().
		//WithLowAccuracyMode().
		Build()

	return detector, nil
}

func DefaultLanguageDetectionOptions() *api.LanguageDetectionOptions {
	return &api.LanguageDetectionOptions{
		DefautLanguageCode:   "en",
		DefautLanguageName:   "english",
		ConfidenceThresthold: 0.20,
		MinimumWordCount:     5,
	}
}

var _ api.LanguageDetector = &languageDetector{}

type languageDetector struct {
	detector            lingua.LanguageDetector
	defaultLanguageCode string
	defaultLanguageName string
	minimumWordCount    int
	threshold           float64
}

func (i *languageDetector) DetectLanguage(sentence string, ignoreWordCount bool) (string, string, bool) {
	simpleWordCount := whitespaceRegexp.Split(sentence, -1)
	log.Debug().Msgf("detect language with %d words (ignoreWordCount=%v)", len(simpleWordCount), ignoreWordCount)

	if len(simpleWordCount) <= i.minimumWordCount && !ignoreWordCount {
		log.Debug().Msgf("word cound %d is below thresthold of %d, return default values", len(simpleWordCount), i.minimumWordCount)

		// below 5 word the language detection becomes bad
		// see https://github.com/pemistahl/lingua-go?tab=readme-ov-file#4-how-good-is-it
		// return detector defaults
		return i.defaultLanguageCode, i.defaultLanguageName, true
	}

	confidenceValues := i.detector.ComputeLanguageConfidenceValues(sentence)
	log.Debug().Msgf(
		"detected language \"%s\" with confidence %f",
		strings.ToLower(confidenceValues[0].Language().IsoCode639_1().String()),
		confidenceValues[0].Value(),
	)

	bigFloatHighest := big.NewFloat(confidenceValues[0].Value())
	bigFloatThreshold := big.NewFloat(i.threshold)
	result := bigFloatThreshold.Cmp(bigFloatHighest)
	log.Debug().Msgf("over the detection thresthold (%f) => %v", i.threshold, result < 0)

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
