package api

type Tokenizer interface {
	AddDictionary(dictionary *Dictionary) error
	Tokenize(sentence string) ([]*Token, error)
}
