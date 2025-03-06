package config

import "github.com/ylallemant/githook-companion/pkg/api"

const (
	typeFeature  = "feat"
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
	commit.Dictionaries = commitDictionaries()
	commit.Types = commitTypes()
	commit.LanguageCodes = []string{
		"en",
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
			Description: "introducing a breaking change",
		},
	}
}

func commitDictionaries() []*api.CommitTypeDictionary {
	return []*api.CommitTypeDictionary{
		{
			Name:  "add",
			Value: "add",
			Type:  typeFeature,
			Synonyms: []string{
				"adds",
				"added",
				"adding",
				"new",
			},
		},
		{
			Name:  "use",
			Value: "use",
			Type:  typeFeature,
			Synonyms: []string{
				"used",
				"uses",
			},
		},
		{
			Name:  "update",
			Value: "update",
			Type:  typeRefactor,
			Synonyms: []string{
				"updated",
				"updates",
			},
		},
		{
			Name:  typeRefactor,
			Value: typeRefactor,
			Type:  typeRefactor,
			Synonyms: []string{
				"change",
				"changes",
				"changed",
				"restructure",
				"restructured",
				"restructures",
			},
		},
		{
			Name:  "remove",
			Value: "remove",
			Type:  typeRefactor,
			Synonyms: []string{
				"removed",
				"removes",
			},
		},
		{
			Name:  "fix",
			Value: "fix",
			Type:  typeFeature,
			Synonyms: []string{
				"fixes",
				"fixed",
				"fixing",
				"correct",
			},
		},
	}
}
