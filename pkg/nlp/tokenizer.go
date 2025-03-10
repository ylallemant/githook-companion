package nlp

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/nlp/api"
)

var _ api.Tokenizer = &tokenizer{}

type tokenizer struct {
	languageDetector     *LanguageDetector
	splitters            map[string]*splitter
	normalisers          map[string]*normaliser
	dictionaries         []*api.Dictionary
	lexemes              []*api.Lexeme
	confidenceThresthold float64
}

func NewTokenizer(options *api.TokenizerOptions) (*tokenizer, error) {
	instance := new(tokenizer)

	if options == nil {
		options = DefaultTokenizerOptions()
	}

	if options.ConfidenceThresthold > 0 {
		instance.confidenceThresthold = options.ConfidenceThresthold
	} else {
		instance.confidenceThresthold = DefaultConfidenceThresthold
	}

	if len(options.LanguageCodes) < 2 {
		options.LanguageCodes = append(options.LanguageCodes, []string{"en", "de"}...)
	}

	languageDetector, err := NewLanguageDetector(options.LanguageCodes, 0.8)
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

func (i *tokenizer) AddDictionary(dictionary *api.Dictionary) error {
	for _, current := range i.dictionaries {
		if current.Name == dictionary.Name {
			return errors.Errorf("failed to add dictionary, \"%s\" already registered", dictionary.Name)
		}
	}

	i.dictionaries = append(i.dictionaries, dictionary)

	return nil
}

func (i *tokenizer) Tokenize(sentence string) ([]*api.Token, string, error) {
	if sentence == "" {
		return []*api.Token{}, "", nil
	}

	lagnCode, _, known := i.languageDetector.DetectLanguage(sentence)

	if !known {
		return []*api.Token{}, "", errors.Errorf("failed to detect language from sentence : \"%s\"", sentence)
	}

	template, words := i.split(sentence, lagnCode)

	i.normalise(words, lagnCode)

	tokens := i.matchTokens(words)

	return tokens, template, nil
}

func (i *tokenizer) matchTokens(words []*api.Word) []*api.Token {
	tokens := make([]*api.Token, 0)

	for _, word := range words {
		token := new(api.Token)
		token.Word = word

		if word.FromLexeme != "" {
			token.Name = word.FromLexeme
			token.Value = word.Normalised
			token.Source = api.TokenSourceLexeme
			token.Confidence = 1
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

func (i *tokenizer) split(sentence, languageCode string) (string, []*api.Word) {
	if splitter, ok := i.splitters[languageCode]; ok {
		return splitter.Split(sentence)
	}

	return i.splitters[api.LanguageCodeWildcard].Split(sentence)
}

func (i *tokenizer) normalise(words []*api.Word, languageCode string) {
	if normaliser, ok := i.normalisers[languageCode]; ok {
		normaliser.NormaliseAll(words)
	} else {
		panic(fmt.Sprintf("unknown language-code \"%s\"", languageCode))
	}
}

// fuzzyDictionaryMatch returns the matching dictionary,
// diectionary entry and confidence score if any.
// A threshold can be set in the tokenizer options.
func (i *tokenizer) fuzzyDictionaryMatch(word *api.Word) (*api.Dictionary, string, float64) {
	var match *api.Dictionary
	bestConfidence := 0.0
	bestMatch := word.Normalised

	for _, dictionary := range i.dictionaries {
		if dictionary.LanguageCode != word.LanguageCode && dictionary.LanguageCode != api.LanguageCodeWildcard {
			// dictionary is not relevant, skip it
			continue
		}

		for _, entry := range dictionary.Entries {
			confidence := calculateConfidence(word.Normalised, entry)

			if confidence > bestConfidence && confidence > dictionary.ConfidenceThresthold {
				bestConfidence = confidence
				match = dictionary
				bestMatch = entry
			}
		}
	}

	if bestConfidence <= i.confidenceThresthold {
		return nil, "", 0.0
	}

	return match, bestMatch, bestConfidence
}
