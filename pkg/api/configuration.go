package api

import "regexp"

const ConfigDirectory = ".githook-companion"
const ConfigFile = "config.yaml"

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
	Dictionaries    []*CommitTypeDictionary `yaml:"dictionaries" json:"dictionaries"`
	Types           []*CommitType           `yaml:"types" json:"types"`
	Tokens          []*CommitToken          `yaml:"tokens" json:"tokens"`
	DefaultType     string                  `yaml:"defaut_type" json:"defaut_type"`
	MessageTemplate string                  `yaml:"message_template" json:"message_template"`
}

type CommitType struct {
	Type        string `yaml:"type" json:"type"`
	Description string `yaml:"description" json:"description"`
}

type CommitTypeDictionary struct {
	Name     string   `yaml:"name" json:"name"`
	Value    string   `yaml:"value" json:"value"`
	Type     string   `yaml:"type" json:"type"`
	Synonyms []string `yaml:"synonyms" json:"synonyms"`
}

type CommitToken struct {
	Name    string
	Lexemes []*regexp.Regexp
}
