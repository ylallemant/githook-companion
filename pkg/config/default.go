package config

import (
	"github.com/ylallemant/githook-companion/pkg/api"
	nlpapi "github.com/ylallemant/githook-companion/pkg/nlp/api"
)

const (
	typeFeature  = "feat"
	typeIgnore   = "ignore"
	typeDocs     = "docs"
	typeFix      = "fix"
	typeTest     = "test"
	typeRefactor = "refactor"
	typeBreaking = "breaking"
)

func Default() *api.Config {
	config := new(api.Config)

	commit := new(api.Commit)
	commit.DefaultType = typeFeature
	commit.Types = commitTypes()
	commit.TokenizerOptions = &nlpapi.TokenizerOptions{
		LanguageCodes: []string{
			"en",
		},
		Dictionaries: commitDictionaries(),
	}

	config.Commit = commit

	return config
}

func commitTypes() []*api.CommitType {
	return []*api.CommitType{
		{
			Type:        typeFeature,
			Description: "a new feature is introduced with the changes",
		},
		{
			Type:        typeIgnore,
			Description: "commit can be ignored by other tools",
		},
		{
			Type:        typeFix,
			Description: "a bug fix has occurred",
		},
		{
			Type:        typeDocs,
			Description: "updates to documentation such as a the README or other markdown files",
		},
		{
			Type:        typeTest,
			Description: "including new or correcting previous tests",
		},
		{
			Type:        typeRefactor,
			Description: "refactored code that neither fixes a bug nor adds a feature",
		},
		{
			Type:        typeBreaking,
			Description: "introducing a breaking change in input or output behaviour",
		},
	}
}

func commitDictionaries() []*nlpapi.Dictionary {
	return []*nlpapi.Dictionary{
		{
			LanguageCode: "en",
			Name:         "weak-feature-signals",
			TokenName:    typeFeature,
			Entries: []string{
				"add",
				"implement",
				"use",
				"new",
			},
		},
		{
			LanguageCode: "en",
			Name:         "ignore-signals",
			TokenName:    typeIgnore,
			Entries: []string{
				"typo",
			},
		},
		{
			LanguageCode: "en",
			Name:         typeRefactor,
			TokenName:    typeRefactor,
			Entries: []string{
				"remove",
				"change",
				"update",
				"upgrate",
				"restructure",
			},
		},
		{
			LanguageCode: "en",
			Name:         "fix",
			TokenName:    typeFix,
			Entries: []string{
				"fix",
			},
		},
	}
}
