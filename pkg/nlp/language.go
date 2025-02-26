package nlp

import "github.com/rylans/getlang"

func DetectLanguage(sentence string) (string, string, float64) {
	i18n := getlang.FromString(sentence)

	return i18n.LanguageCode(), i18n.LanguageName(), i18n.Confidence()
}
