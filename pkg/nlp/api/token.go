package api

const (
	TokenUnknown          = "unknown"
	TokenSourceNone       = "none"
	TokenSourceLexeme     = "lexeme"
	TokenSourceDictionary = "dictionary"
)

type Token struct {
	Word       *Word
	Name       string
	Source     string
	Value      string
	Confidence float64
}
