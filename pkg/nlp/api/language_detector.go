package api

type LanguageDetector interface {
	DetectLanguage(sentence string, strict bool) (string, string, bool)
}
