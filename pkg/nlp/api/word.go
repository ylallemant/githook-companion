package api

const (
	WordSourceLexeme         = "lexeme"
	WordSourceLexemeSplitter = "lexeme-splitter"
	WordSourceSplitter       = "splitter"
)

type Word struct {
	LanguageCode string
	Raw          string
	Cleaned      string
	Source       string
	SourceName   string
	Normalised   string
}
