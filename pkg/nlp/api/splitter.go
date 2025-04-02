package api

import "regexp"

const (
	LexemeReferenceFmt = "lexeme~%d~%d"
	WordReferenceFmt   = "word~%d"
)

var (
	LexemeKeyRegexp   = regexp.MustCompile(`lexeme~\d+~\d+`)
	WordKeyRegexp     = regexp.MustCompile(`word~\d+`)
	PlaceholderRegexp = regexp.MustCompile(`((lexeme|word)~\d+(~\d+){0,1})`)
)

type Splitter interface {
	LanguageCode() string
	Split(sentence string) (string, []*Word)
}
