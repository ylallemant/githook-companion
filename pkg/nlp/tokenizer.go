package nlp

import (
	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/nlp/api"
)

var _ api.Tokenizer = &tokanizer{}

type tokanizer struct {
	dictionaries []*api.Dictionary
}

func Tokinizer() *tokanizer {
	return new(tokanizer)
}

func (i *tokanizer) Tokenize(sentence string) ([]*api.Token, error) {
	tokens := make([]*api.Token, 0)

	lagnCode, _, known := DetectLanguage(sentence)

	if !known {
		return tokens, errors.Errorf("failed to detect language from sentence : \"%s\"", sentence)
	}

	_, err := Lemmatizer(lagnCode)
	if err != nil {
		return tokens, err
	}

	return tokens, nil
}

func (i *tokanizer) AddDictionary(dictionary *api.Dictionary) error {
	for _, current := range i.dictionaries {
		if current.Name == dictionary.Name {
			return errors.Errorf("failed to add dictionary, \"%s\" already registered", dictionary.Name)
		}
	}

	i.dictionaries = append(i.dictionaries, dictionary)

	return nil
}
