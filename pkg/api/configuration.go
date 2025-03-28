package api

import (
	"github.com/pkg/errors"
	nlpapi "github.com/ylallemant/githook-companion/pkg/nlp/api"
)

const (
	ConfigDirectory     = ".githook-companion"
	ConfigFile          = "config.yaml"
	GithooksDirectory   = "hooks"
	BinDirectory        = "bin"
	ContextDirectory    = "context"
	LogDirectory        = "logs"
	CommitTypeTokenName = "commit_type"
	CommitMessageKey    = "message"
)

var (
	ConfigurationNotFound       = errors.New("configuration not found")
	ConfigProcessingDirectories = []string{
		LogDirectory,
		BinDirectory,
		ContextDirectory,
	}
)

type Config struct {
	*ParentConfig       `yaml:"parent,omitempty" json:"parent,omitempty"`
	*Commit             `yaml:"commit" json:"commit"`
	Dependencies        []*Dependency `yaml:"dependencies" json:"dependencies"`
	DependencyDirectory string        `yaml:"dependency_directory" json:"dependency_directory"`
	GithooksDirectory   string        `yaml:"githook_directory" json:"githook_directory"`
}

type ParentConfig struct {
	Path          string `yaml:"path" json:"path"`
	GitRepository string `yaml:"repository" json:"repository"`
	Private       bool   `yaml:"private" json:"private"`
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
	NoFormatting     []string                 `yaml:"no_formatting" json:"no_formatting"`
	DefaultType      string                   `yaml:"default_type" json:"default_type"`
	MessageTemplate  string                   `yaml:"message_template" json:"message_template"`
	TokenizerOptions *nlpapi.TokenizerOptions `yaml:"tokenizer_options" json:"tokenizer_options"`
}

type CommitType struct {
	Type        string `yaml:"type" json:"type"`
	Description string `yaml:"description" json:"description"`
}
