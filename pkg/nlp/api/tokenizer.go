package api

var LanguageCodeWildcard = "*"

type Tokenizer interface {
	AddDictionary(dictionary *Dictionary) error
	Tokenize(sentence string) ([]*Token, string, error)
}

type TokenizerOptions struct {
	LanguageCodes        []string
	Dictionaries         []*Dictionary
	Lexemes              []*Lexeme
	Normalisers          []*Normaliser
	ConfidenceThresthold float64
}
