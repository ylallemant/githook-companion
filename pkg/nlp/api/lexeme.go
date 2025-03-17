package api

import (
	"regexp"
)

var DefaultMatcherRegexp = regexp.MustCompile(".*")

type Lexeme struct {
	Name         string               `yaml:"name" json:"name"`
	Description  string               `yaml:"description" json:"description"`
	LanguageCode string               `yaml:"language_code" json:"language_code"`
	TokenName    string               `yaml:"token" json:"token"`
	Variants     []*Variant           `yaml:"variants" json:"variants"`
	Normalisers  []*NormalisationStep `yaml:"normalisers" json:"normalisers"`
}

type Variant struct {
	Name        string               `yaml:"name" json:"name"`
	Description string               `yaml:"description" json:"description"`
	Matcher     *Matcher             `yaml:"matcher" json:"matcher"`
	Normalisers []*NormalisationStep `yaml:"normalisers" json:"normalisers"`
}

type Matcher struct {
	Regex *regexp.Regexp
}

// UnmarshalText unmarshals json into a regexp.Regexp
func (r *Matcher) UnmarshalText(b []byte) error {
	regex := DefaultMatcherRegexp
	var err error

	if len(b) > 0 {
		regex, err = regexp.Compile(string(b))
		if err != nil {
			return err
		}
	}

	r.Regex = regex

	return nil
}

// MarshalText marshals regexp.Regexp as string
func (r *Matcher) MarshalText() ([]byte, error) {
	if r.Regex != nil {
		return []byte(r.Regex.String()), nil
	}

	return nil, nil
}
