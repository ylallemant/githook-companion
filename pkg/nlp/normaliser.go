package nlp

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/nlp/api"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	asciiRegexp   = regexp.MustCompile(`[^a-z0-9\s]`)
	acronymRegexp = regexp.MustCompile(`\b(?:[a-zA-Z]\.){2,}([a-zA-Z]){0,1}`)
)

func DefaultNormaliser(languageCode string) (*normaliser, error) {
	instance := new(normaliser)

	instance.languageCode = languageCode

	var err error
	instance.lemmatizer, err = Lemmatizer(languageCode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to instantiate default normaliser")
	}

	return instance, nil
}

var _ api.Normaliser = &normaliser{}

type normaliser struct {
	languageCode string
	lemmatizer   api.Lemmatizer
}

func (i *normaliser) LanguageCode() string {
	return i.languageCode
}

func (i *normaliser) NormaliseAll(words []*api.Word) {
	for _, word := range words {
		i.Normalise(word)
	}
}

func (i *normaliser) Clean(word *api.Word) {
	word.Cleaned = strings.ToLower(word.Raw)
	word.Cleaned = removeDiacritics(word.Cleaned)
	word.Cleaned = replaceAcronyms(word.Cleaned)
	word.Cleaned = removeSpecialChars(word.Cleaned)
	word.Cleaned = strings.TrimSpace(word.Cleaned)
}

func (i *normaliser) Normalise(word *api.Word) {
	if word.FromLexeme != "" {
		return
	}

	i.Clean(word)
	i.lemmatizer.Lemma(word)
}

func removeDiacritics(text string) string {
	transformer := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(transformer, text)
	return result
}

func replaceAcronyms(text string) string {
	acronyms := acronymRegexp.FindAll([]byte(text), -1)

	for _, acronym := range acronyms {
		cleaned := strings.ReplaceAll(string(acronym), ".", "")
		text = strings.ReplaceAll(text, string(acronym), cleaned)
	}

	return text
}

func removeSpecialChars(text string) string {
	text = asciiRegexp.ReplaceAllString(text, " ")
	text = whitespaceRegexp.ReplaceAllString(text, " ")

	return text
}

// addToMap - add src's entries to dst
func addToMap(dst, src map[string]interface{}) {
	for k, v := range src {
		dst[k] = v
	}
}
