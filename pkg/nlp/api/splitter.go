package api

import "regexp"

const (
	LexemeReferenceFmt = "lexeme:%d"
	WordReferenceFmt   = "word:%d"
)

var (
	LexemeKeyRegexp = regexp.MustCompile(`lexeme:\d+`)
)

type Splitter interface {
	LanguageCode() string
	Split(sentence string) (string, []*Word)
}
