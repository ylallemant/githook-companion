package api

import (
	"crypto/sha256"
	"encoding/hex"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/pkg/errors"
)

type Normaliser interface {
	LanguageCode() string
	NormaliseAll(words []*Word)
	Normalise(word *Word)
	Clean(word *Word)
}

// NormalisationStep is used to modify the provided text.
// the modification affects the text in parts or completly
type NormalisationStep struct {
	// Name of the normalisation step
	Name string `yaml:"name" json:"name"`
	// Matcher a regular-exception to select the part[s]
	// of the text subjected to the normalisation step
	// it returns a list of strings
	Matcher *Matcher `yaml:"matcher" json:"matcher"`
	// ReplaceAll specifies if the whole text has to be replaced
	// by the first part returned by the Matcher
	ReplaceAll bool `yaml:"selective" json:"selective"`
	// Replacement will replace all parts returned by Matcher
	Replacement string `yaml:"replacement" json:"replacement"`
	// Formatter allows the usage of go-templates to format the selected parts
	Formatter *Formatter `yaml:"formatter" json:"formatter"`
}

// Formatter makes use of go templates to modify the provided text.
// build-in template functions by Masterminds Sprig package
// documentation can be found here https://masterminds.github.io/sprig/
type Formatter struct {
	Template string `yaml:"template" json:"template"`
	Renderer *template.Template
}

// UnmarshalText unmarshals json into a regexp.Regexp
func (r *Formatter) UnmarshalText(b []byte) error {
	if r.Template != "" {
		b = []byte(r.Template)

		hasher := sha256.New()
		hasher.Write(b)

		templateName := hex.EncodeToString(hasher.Sum(nil))

		tmpl := template.New(templateName).Funcs(sprig.FuncMap())

		var err error
		if tmpl, err = tmpl.Parse(string(b)); err != nil {
			return errors.Wrapf(err, "failed to parse template \"%s\"", string(b))
		}

		r.Renderer = tmpl
	}

	return nil
}

// MarshalText marshals regexp.Regexp as string
func (r *Formatter) MarshalText() ([]byte, error) {
	return nil, nil
}
