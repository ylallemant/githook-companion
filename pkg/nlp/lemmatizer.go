package nlp

import (
	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/de"
	"github.com/aaaton/golem/v4/dicts/en"
	"github.com/aaaton/golem/v4/dicts/fr"
	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/nlp/api"
)

var _ api.Lemmatizer = &lemmatizer{}

func Lemmatizer(i18n string) (api.Lemmatizer, error) {
	l := new(lemmatizer)
	var langpack golem.LanguagePack
	var err error

	switch i18n {
	case "en":
		langpack = en.New()
	case "fr":
		langpack = fr.New()
	case "de":
		langpack = de.New()
	default:
		return nil, errors.Errorf("unsupported language \"%s\"", i18n)
	}

	l.tool, err = golem.New(langpack)
	if err != nil {
		return nil, errors.Errorf("fialed to instantiate lemmatizer for language \"%s\"", i18n)
	}

	return l, nil
}

type lemmatizer struct {
	tool *golem.Lemmatizer
}

func (l *lemmatizer) Lemma(word string) string {
	return l.tool.Lemma(word)
}
