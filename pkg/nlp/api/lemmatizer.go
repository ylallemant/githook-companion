package api

type Lemmatizer interface {
	Lemma(word string) string
}
