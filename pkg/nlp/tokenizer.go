package nlp

import (
	"github.com/pkg/errors"
	"github.com/pkg/nlp/api"
)

func Tokenize(sentence string) ([]*api.Token, error) {
	tokens := make([]*api.Token, 0)

	lagnCode, _, confidence := DetectLanguage(sentence)

	if lagnCode == "" || confidence < 0 {
		return tokens, errors.Errorf("failed to detect language from sentence : \"%s\"", sentence)
	}

	lemmatizer, err := New(lagnCode)
	if err != nil {
		return tokens, err
	}

	return tokens, nil
}
