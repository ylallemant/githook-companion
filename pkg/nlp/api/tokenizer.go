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
	LanguageDetector() LanguageDetector
}

type TokenizerOptions struct {
	LanguageDetectionOptions *LanguageDetectionOptions `yaml:"language_detection_options" json:"language_detection_options"`
	LanguageCodes            []string                  `yaml:"language_codes" json:"language_codes"`
	Dictionaries             []*Dictionary             `yaml:"dictionaries" json:"dictionaries"`
	Lexemes                  []*Lexeme                 `yaml:"lexemes" json:"lexemes"`
	Normalisers              []*Normaliser             `yaml:"normalisers" json:"normalisers"`
	ConfidenceThresthold     float64                   `yaml:"confidence_thresthold" json:"confidence_thresthold"`
}

type LanguageDetectionOptions struct {
	DefautLanguageCode   string  `yaml:"default_language_code" json:"default_language_code"`
	DefautLanguageName   string  `yaml:"default_language_Name" json:"default_language_Name"`
	ConfidenceThresthold float64 `yaml:"confidence_thresthold" json:"confidence_thresthold"`
	MinimumWordCount     int     `yaml:"minimum_word_count" json:"minimum_word_count"`
}
