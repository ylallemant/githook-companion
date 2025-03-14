package api

import "regexp"

var (
	LanguageCodeWildcard = "*"
	TokenNameRegexp      = regexp.MustCompile("[a-zA-Z][a-zA-Z0-9_]+")
)

type Tokenizer interface {
	AddLexeme(lexeme *Lexeme) error
	AddDictionary(dictionary *Dictionary) error
	Tokenize(sentence string) ([]*Token, string, string, error)
	ValidateTokenName(name string) bool
}

type TokenizerOptions struct {
	LanguageCodes        []string
	Dictionaries         []*Dictionary
	Lexemes              []*Lexeme
	Normalisers          []*Normaliser
	ConfidenceThresthold float64
}
