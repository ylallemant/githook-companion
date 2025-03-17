package api

import "regexp"

const (
	LexemeReferenceFmt = "lexeme~%d"
	WordReferenceFmt   = "word~%d"
)

var (
	LexemeKeyRegexp   = regexp.MustCompile(`lexeme~\d+`)
	WordKeyRegexp     = regexp.MustCompile(`word~\d+`)
	PlaceholderRegexp = regexp.MustCompile(`((lexeme|word)~\d+)`)
)

type Splitter interface {
	LanguageCode() string
	Split(sentence string) (string, []*Word)
}
