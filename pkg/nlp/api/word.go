package api

const (
	WordSourceLexeme   = "lexeme"
	WordSourceSplitter = "splitter"
)

type Word struct {
	LanguageCode string
	Raw          string
	Cleaned      string
	Source       string
	SourceName   string
	Normalised   string
}
