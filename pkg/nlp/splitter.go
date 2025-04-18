package nlp

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/ylallemant/githook-companion/pkg/nlp/api"
	nlpapi "github.com/ylallemant/githook-companion/pkg/nlp/api"
)

var (
	nonScriptCharRegexp = regexp.MustCompile(`[^\p{L}]`)
	whitespaceRegexp    = regexp.MustCompile(`\s+`)
	puntuationRegexp    = regexp.MustCompile(`([^0-9\p{L}\s_-~])`)
	minWordLength       = 3
)

func DefaultSplitter(languageCode string, lexemes []*api.Lexeme) *splitter {
	instance := new(splitter)

	instance.minWordLength = minWordLength
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
	languageCode  string
	minWordLength int
	lexemes       []*api.Lexeme
}

func (i *splitter) LanguageCode() string {
	return i.languageCode
}

func (i *splitter) Split(sentence string) (string, []*api.Word) {
	sentence = strings.TrimSpace(sentence)

	// extract complex lexemes and replace them with position information
	sentenceTemplate, wordsFromLexemes := i.ExtractLexemes(sentence)
	log.Debug().Msgf("found %d lexemes in sentence", len(wordsFromLexemes))

	// clean sentence without messing with diacritics
	splitTemplate := i.clean(sentenceTemplate)
	log.Debug().Msgf("cleaned sentence: %s", splitTemplate)

	words := make([]*api.Word, 0)

	parts := whitespaceRegexp.Split(splitTemplate, -1)
	log.Debug().Msgf("splitted sentence: %v", parts)

	for _, part := range parts {
		if len(part) < i.minWordLength {
			// ignore sigle characters
			log.Debug().Msgf("ignore single character: \"%s\"", part)
			continue
		}

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

	log.Debug().Msgf("check sentence for %d lexemes", len(i.lexemes))
	for _, lexeme := range i.lexemes {
		for _, matcher := range lexeme.Variants {
			log.Debug().Msgf("search for lexeme \"%s\" variant with: \"%s\"", lexeme.Name, matcher.Matcher.Regex.String())
			if matcher.Matcher.Regex.MatchString(sentence) {
				matches := FindAll(sentence, matcher.Matcher)
				log.Debug().Msgf("  - matches: %s", matches)

				for _, match := range matches {
					word := new(api.Word)

					word.LanguageCode = lexeme.LanguageCode
					word.Raw = strings.TrimSpace(match)
					word.Source = api.WordSourceLexeme
					word.SourceName = lexeme.TokenName

					i.normaliseLexeme(word, matcher.Normalisers, lexeme.Normalisers)

					replacementWords := make(map[string]*api.Word)

					if len(lexeme.Splitters) > 0 {
						allwords, _ := i.splitLexeme(word.Normalised, len(words), lexeme)

						for key, word := range allwords {
							// add parent lexeme match
							word.FromComposite = strings.TrimSpace(match)

							words[key] = word
							replacementWords[key] = word
						}
					} else {
						key := fmt.Sprintf(api.LexemeReferenceFmt, len(words), 0)

						words[key] = word
						replacementWords[key] = word
					}

					for key, word := range replacementWords {
						// replace lexeme with index information
						// add spaces as prefix and suffix to make sure
						// the splitter will be able to split
						log.Debug().Msgf("    - replacement: %s, %s", word.Raw, key)
						sentence = secureReplaceAllString(sentence, word.Raw, fmt.Sprintf(" %s ", key))
						log.Debug().Msgf("    - sentence: %s", sentence)
					}
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

func (i *splitter) normaliseLexeme(word *api.Word, primary []*api.NormalisationStep, secondary []*api.NormalisationStep) {
	text := word.Raw

	for _, normalisationStep := range primary {
		text = executeNormaliser(normalisationStep, text)
	}

	for _, normalisationStep := range secondary {
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
				err := normaliser.Formatter.UnmarshalText([]byte{})
				if err != nil {
					panic(fmt.Sprintf(
						"failed to unmarshal formatter with template \"%s\": %s",
						normaliser.Formatter.Template,
						err.Error(),
					))
				}
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

func (i *splitter) splitLexeme(text string, lexemeIndex int, lexeme *nlpapi.Lexeme) (map[string]*api.Word, error) {
	words := make(map[string]*api.Word)
	template := text

	for _, splitter := range lexeme.Splitters {

		for _, variant := range splitter.Variants {
			log.Debug().Msgf("    - search for part \"%s\" variant with: \"%s\"", splitter.Name, variant.Matcher.Regex.String())

			if Match(text, variant.Matcher) {
				matches := FindAll(text, variant.Matcher)
				log.Debug().Msgf("      - matches: %s", matches)

				for _, match := range matches {
					if match == "" {
						continue
					}

					log.Debug().Msgf("        - match: %s", match)
					word := new(api.Word)

					word.LanguageCode = lexeme.LanguageCode
					word.Raw = strings.TrimSpace(match)
					word.Source = api.WordSourceLexemeSplitter
					word.SourceName = splitter.TokenName
					log.Debug().Msgf("        - raw: %s", word.Raw)

					key := fmt.Sprintf(api.LexemeReferenceFmt, lexemeIndex, len(words))
					log.Debug().Msgf("        - key: %s", key)

					i.normaliseLexeme(word, variant.Normalisers, splitter.Normalisers)
					log.Debug().Msgf("        - nomalised: %s", word.Normalised)

					words[key] = word

					// replace lexeme with index information
					// add spaces as prefix and suffix to make sure
					// the splitter will be able to split
					template = secureReplaceAllString(template, word.Raw, fmt.Sprintf(" %s ", key))
				}
			}
		}
	}

	return words, nil
}
