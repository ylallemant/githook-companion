package api

type Word struct {
	LanguageCode string
	Raw          string
	Cleaned      string
	FromLexeme   string
	Normalised   string
}
