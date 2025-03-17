package api

import (
	nlpapi "github.com/ylallemant/githook-companion/pkg/nlp/api"
)

const (
	ConfigDirectory     = ".githook-companion"
	ConfigFile          = "config.yaml"
	CommitTypeTokenName = "commit_type"
	CommitMessageKey    = "message"
)

type Config struct {
	*Commit      `yaml:"commit" json:"commit"`
	Dependencies []*Dependency `yaml:"dependencies" json:"dependencies"`
}

type Dependency struct {
	Name                string   `yaml:"name" json:"name"`
	Version             string   `yaml:"version" json:"version"`
	SemverPrefix        string   `yaml:"semver_prefix" json:"semver_prefix"`
	Args                []string `yaml:"args" json:"args"`
	ForceReplace        bool     `yaml:"force-replace" json:"force-replace"`
	ChecksumUrlTemplate string   `yaml:"checksum_url_template" json:"checksum_url_template"`
	ChecksumAlgorithm   string   `yaml:"checksum_algorithm" json:"checksum_algorithm"`
	ArchiveUrlTemplate  string   `yaml:"archive_url_template" json:"archive_url_template"`
}

type Commit struct {
	Types            []*CommitType            `yaml:"types" json:"types"`
	DefaultType      string                   `yaml:"defaut_type" json:"defaut_type"`
	MessageTemplate  string                   `yaml:"message_template" json:"message_template"`
	TokenizerOptions *nlpapi.TokenizerOptions `yaml:"tokenizer_options" json:"tokenizer_options"`
}

type CommitType struct {
	Type        string `yaml:"type" json:"type"`
	Description string `yaml:"description" json:"description"`
}
