package api

const (
	TokenUnknown               = "unknown"
	TokenSourceNone            = "none"
	TokenSourceLexeme          = "lexeme"
	TokenSourceLexemeComposite = "lexeme-composite"
	TokenSourceDictionary      = "dictionary"
)

type Token struct {
	Word        *Word
	Name        string
	Source      string
	SourceName  string
	SourceMatch string
	Value       string
	Confidence  float64
	Index       int64
}
