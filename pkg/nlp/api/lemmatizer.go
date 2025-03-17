package api

type Lemmatizer interface {
	LanguageCode() string
	Lemma(word *Word)
}
