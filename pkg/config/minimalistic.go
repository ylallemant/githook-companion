package config

import (
	"github.com/ylallemant/githook-companion/pkg/api"
	nlpapi "github.com/ylallemant/githook-companion/pkg/nlp/api"
)

func Minimalistic() *api.Config {
	config := new(api.Config)
	config.Commit = new(api.Commit)
	config.Commit.TokenizerOptions = new(nlpapi.TokenizerOptions)

	return config
}
