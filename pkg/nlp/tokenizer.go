package nlp

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/nlp/api"
)

var _ api.Tokenizer = &tokanizer{}

type tokanizer struct {
	languageDetector *LanguageDetector
	splitters        map[string]*splitter
	normalisers      map[string]*normaliser
	dictionaries     []*api.Dictionary
	lexemes          []*api.Lexeme
}

func Tokinizer(options *api.TokenizerOptions) (*tokanizer, error) {
	instance := new(tokanizer)

	if len(options.LanguageCodes) == 0 {
		options.LanguageCodes = []string{"en", "de"}
	}

	languageDetector, err := NewLanguageDetector(options.LanguageCodes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initiate the Tokenizer")
	}

	instance.languageDetector = languageDetector

	for _, dictionary := range options.Dictionaries {
		err = instance.AddDictionary(dictionary)
		if err != nil {
			return nil, errors.Wrap(err, "failed to import dictionaries")
		}
	}

	instance.splitters = make(map[string]*splitter)
	instance.normalisers = make(map[string]*normaliser)

	for _, languageCode := range options.LanguageCodes {
		languageNormaliser, err := DefaultNormaliser(languageCode)
		if err != nil {
			return nil, err
		}

		instance.normalisers[languageCode] = languageNormaliser

		instance.splitters[languageCode] = DefaultSplitter(languageCode, instance.lexemes)
	}

	return instance, nil
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

func (i *tokanizer) Tokenize(sentence string) ([]*api.Token, string, error) {
	lagnCode, _, known := i.languageDetector.DetectLanguage(sentence)

	if !known {
		return []*api.Token{}, "", errors.Errorf("failed to detect language from sentence : \"%s\"", sentence)
	}

	template, words := i.split(sentence, lagnCode)

	i.normalise(words, lagnCode)

	tokens := i.matchTokens(words)

	return tokens, template, nil
}

func (i *tokanizer) matchTokens(words []*api.Word) []*api.Token {
	tokens := make([]*api.Token, 0)

	for _, word := range words {
		token := new(api.Token)
		token.Word = word

		if word.FromLexeme != "" {
			token.Name = word.FromLexeme
			token.Value = word.Normalised
			token.Source = api.TokenSourceLexeme
		} else {
			dictionary, match, confidence := i.fuzzyDictionaryMatch(word)

			token.Confidence = confidence

			if dictionary == nil {
				token.Name = api.TokenUnknown
				token.Value = word.Normalised
				token.Source = api.TokenSourceNone
			} else {
				token.Name = dictionary.TokenName
				token.Value = match
				token.Source = api.TokenSourceDictionary
			}
		}

		tokens = append(tokens, token)
	}

	return tokens
}

func (i *tokanizer) split(sentence, languageCode string) (string, []*api.Word) {
	if splitter, ok := i.splitters[languageCode]; ok {
		return splitter.Split(sentence)
	}

	return i.splitters[api.LanguageCodeWildcard].Split(sentence)
}

func (i *tokanizer) normalise(words []*api.Word, languageCode string) {
	if normaliser, ok := i.normalisers[languageCode]; ok {
		normaliser.NormaliseAll(words)
	}

	panic(fmt.Sprintf("unknown language-code \"%s\"", languageCode))
}

func (i *tokanizer) fuzzyDictionaryMatch(word *api.Word) (*api.Dictionary, string, float64) {
	var match *api.Dictionary
	minDistance := 100000
	minConfidence := 0.0
	bestMatch := word.Normalised

	fmt.Println("dictionary count", len(i.dictionaries))

	for _, dictionary := range i.dictionaries {
		if dictionary.LanguageCode != word.LanguageCode && dictionary.LanguageCode != api.LanguageCodeWildcard {
			// dictionary is not relevant, skip it
			fmt.Println("  - skipping dictionary", dictionary.LanguageCode)
			continue
		}

		for _, entry := range dictionary.Entries {
			confidence := calculateConfidence(word.Normalised, entry)

			fmt.Println("   -", dictionary.Name, word.Normalised, entry, confidence)
			if confidence > minConfidence {
				minConfidence = confidence
				match = dictionary
			}
		}
	}

	if minDistance > 2 {
		return nil, "", 0.0
	}

	return match, bestMatch, 0
}
