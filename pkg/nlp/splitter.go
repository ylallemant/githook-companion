package nlp

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/ylallemant/githook-companion/pkg/nlp/api"
)

var (
	nonScriptCharRegexp = regexp.MustCompile(`[^\p{L}]`)
	whitespaceRegexp    = regexp.MustCompile(`\s+`)
	puntuationRegexp    = regexp.MustCompile(`([^0-9\p{L}\s_-~])`)
)

func DefaultSplitter(languageCode string, lexemes []*api.Lexeme) *splitter {
	instance := new(splitter)

	instance.languageCode = languageCode
	instance.lexemes = make([]*api.Lexeme, 0)

	for _, lexeme := range lexemes {
		if lexeme.LanguageCode == languageCode || lexeme.LanguageCode == api.LanguageCodeWildcard {
			instance.lexemes = append(instance.lexemes, lexeme)
		}
	}

	return instance
}

var _ api.Splitter = &splitter{}

type splitter struct {
	languageCode string
	lexemes      []*api.Lexeme
}

func (i *splitter) LanguageCode() string {
	return i.languageCode
}

func (i *splitter) Split(sentence string) (string, []*api.Word) {
	sentence = strings.TrimSpace(sentence)

	// extract complex lexemes and replace them with position information
	sentenceTemplate, wordsFromLexemes := i.ExtractLexemes(sentence)

	// clean sentence without messing with diacritics
	splitTemplate := i.clean(sentenceTemplate)

	words := make([]*api.Word, 0)

	parts := whitespaceRegexp.Split(splitTemplate, -1)
	for _, part := range parts {
		var word *api.Word
		var rawWord string

		if api.LexemeKeyRegexp.MatchString(part) {
			// keep words from lexemes at the same position in the sentence
			// check if sentence part is a lexeme index reference
			word = wordsFromLexemes[part]
			rawWord = part
		} else {
			// add new word
			word = new(api.Word)

			word.LanguageCode = i.languageCode
			word.Raw = part
			word.Source = api.WordSourceSplitter
			word.SourceName = api.WordSourceSplitter

			rawWord = word.Raw
		}

		// replace raw word by key in the sentence template
		key := fmt.Sprintf(api.WordReferenceFmt, len(words))
		sentenceTemplate = secureReplaceAllString(sentenceTemplate, rawWord, key)

		words = append(words, word)
	}

	return sentenceTemplate, words
}

func (i *splitter) ExtractLexemes(sentence string) (string, map[string]*api.Word) {
	words := make(map[string]*api.Word)

	for _, lexeme := range i.lexemes {
		for _, matcher := range lexeme.Variants {
			if matcher.Matcher.Regex.MatchString(sentence) {
				matches := matcher.Matcher.Regex.FindAllString(sentence, -1)

				for _, match := range matches {
					word := new(api.Word)

					word.LanguageCode = lexeme.LanguageCode
					word.Raw = strings.TrimSpace(match)
					word.Source = api.WordSourceLexeme
					word.SourceName = lexeme.TokenName

					key := fmt.Sprintf(api.LexemeReferenceFmt, len(words))

					i.normaliseLexeme(word, matcher, lexeme)

					words[key] = word

					// replace lexeme with index information
					// add spaces as prefix and suffix to make sure
					// the splitter will be able to split
					sentence = secureReplaceAllString(sentence, word.Raw, fmt.Sprintf(" %s ", key))
				}
			}
		}
	}

	sentence = whitespaceRegexp.ReplaceAllString(sentence, " ")
	sentence = strings.TrimSpace(sentence)

	return sentence, words
}

func (i *splitter) clean(sentence string) string {
	sentence = puntuationRegexp.ReplaceAllString(sentence, " ")
	sentence = whitespaceRegexp.ReplaceAllString(sentence, " ")
	return strings.TrimSpace(sentence)
}

func (i *splitter) normaliseLexeme(word *api.Word, matcher *api.Variant, lexeme *api.Lexeme) {
	text := word.Raw

	for _, normalisationStep := range matcher.Normalisers {
		text = executeNormaliser(normalisationStep, text)
	}

	for _, normalisationStep := range lexeme.Normalisers {
		text = executeNormaliser(normalisationStep, text)
	}

	text = whitespaceRegexp.ReplaceAllString(text, "")

	word.Normalised = strings.TrimSpace(text)
	word.Cleaned = word.Normalised
}

func executeNormaliser(normaliser *api.NormalisationStep, text string) string {
	matches := normaliser.Matcher.Regex.FindAllString(text, -1)

	if normaliser.Replacement != "" {
		text = normaliser.Matcher.Regex.ReplaceAllString(text, normaliser.Replacement)
	}

	if len(matches) > 0 {
		if normaliser.ReplaceAll {
			text = matches[0]
		}

		if normaliser.Formatter != nil {
			if normaliser.Formatter.Renderer == nil {
				normaliser.Formatter.UnmarshalText([]byte{})
			}

			for _, match := range matches {
				var formatted bytes.Buffer
				normaliser.Formatter.Renderer.Execute(&formatted, match)

				if normaliser.ReplaceAll {
					text = formatted.String()
				} else {
					text = strings.ReplaceAll(text, match, formatted.String())
				}
			}
		}
	}

	return text
}
