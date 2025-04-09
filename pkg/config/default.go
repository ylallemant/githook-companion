package config

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/nlp"
	nlpapi "github.com/ylallemant/githook-companion/pkg/nlp/api"
)

const (
	typeFeature     = "feat"
	TypeIgnore      = "ignore"
	typeDocs        = "docs"
	typeFix         = "fix"
	typeTest        = "test"
	typeRefactor    = "refactor"
	typePerformance = "perf"
	typeStyle       = "style"
	typeChore       = "chore"
	typeBuild       = "build"
	typeCI          = "ci"
)

func Default() *api.Config {
	config := new(api.Config)

	config.ToolSync = defaultSync()

	commit := new(api.Commit)
	config.Commit = commit

	commit.MessageTemplate = "{{ .CommitType | upper }}{{ if .CommitScope }}({{ .CommitScope | lower }}){{ end }}{{ if .CommitBreakingFlag }}{{ .CommitBreakingFlag }}{{ end }}: {{ if .IssueTrackerReference }}[{{ .IssueTrackerReference }}]{{ end }} {{ .Message | lower }}"
	commit.DefaultType = typeFeature
	commit.Types = commitTypes()
	commit.NoFormatting = []string{
		TypeIgnore,
	}

	commit.TokenizerOptions = &nlpapi.TokenizerOptions{
		LanguageDetectionOptions: nlp.DefaultLanguageDetectionOptions(),
		LanguageCodes: []string{
			"en",
		},
		Dictionaries: commitDictionaries(),
		Lexemes:      commitLexemes(),
	}

	return config
}

func defaultSync() *api.ToolSync {
	sync := new(api.ToolSync)

	sync.Timeout = 5 * time.Second
	sync.OfflineLockDuration = 20 * time.Minute
	sync.CheckPeriod = 3 * time.Minute

	return sync
}

func commitTypes() []*api.CommitType {
	return []*api.CommitType{
		{
			Type:        typeFeature,
			Description: "a new feature is introduced with the changes",
		},
		{
			Type:        typeRefactor,
			Description: "refactored code that neither fixes a bug nor adds a feature",
		},
		{
			Type:        TypeIgnore,
			Description: "commit can be ignored by other tools",
		},
		{
			Type:        typeFix,
			Description: "a bug fix has been implemented",
		},
		{
			Type:        typeDocs,
			Description: "documentation only changes",
		},
		{
			Type:        typeTest,
			Description: "including new or correcting previous tests",
		},
		{
			Type:        typePerformance,
			Description: "a code change that improves performance",
		},
		{
			Type:        typeStyle,
			Description: "changes that do not affect the meaning of the code (white-space, formatting, ...)",
		},
		{
			Type:        typeChore,
			Description: "other changes that don't modify src or test files",
		},
		{
			Type:        typeBuild,
			Description: "changes that affect the build system or external dependencies",
		},
		{
			Type:        typeCI,
			Description: "changes to CI configuration files and scripts",
		},
	}
}

func CommitTypeList(cfg *api.Config) []string {
	types := make([]string, 0)
	for _, current := range cfg.Commit.Types {
		types = append(types, current.Type)
	}
	return types
}

func commitLexemes() []*nlpapi.Lexeme {
	typeNames := make([]string, 0)

	for _, commitType := range commitTypes() {
		typeNames = append(typeNames, commitType.Type)
	}

	expression := fmt.Sprintf(
		"^(?i)(%s)",
		strings.Join(typeNames, "|"),
	)

	commitTypeReplaceExpression, _ := regexp.Compile(expression)
	commitTypeExpression, _ := regexp.Compile(fmt.Sprintf(
		"%s\\b(\\s*\\(\\w+\\)){0,1}\\s*!{0,1}\\s*:{0,1}",
		expression,
	))

	return []*nlpapi.Lexeme{
		{
			LanguageCode: nlpapi.LanguageCodeWildcard,
			Name:         api.CommitTypeTokenName,
			Description:  "auto-generated commit type lexeme to be retrieved from well formatted messages",
			TokenName:    api.CommitTypeTokenName,
			Variants: []*nlpapi.Variant{
				{
					Matcher: &nlpapi.Matcher{Regex: commitTypeExpression},
				},
			},
			Splitters: []*nlpapi.LexemeSplitter{
				{
					Name:        api.CommitTypeTokenName,
					TokenName:   api.CommitTypeTokenName,
					Description: "commit type lexeme",
					Variants: []*nlpapi.Variant{
						{
							Matcher: &nlpapi.Matcher{Regex: commitTypeReplaceExpression},
						},
					},
					Normalisers: []*nlpapi.NormalisationStep{
						{
							Matcher:    &nlpapi.Matcher{Regex: commitTypeReplaceExpression},
							ReplaceAll: true,
							Formatter: &nlpapi.Formatter{
								Template: "{{ upper . }}",
							},
						},
					},
				},
				{
					Name:        api.CommitScopeTokenName,
					TokenName:   api.CommitScopeTokenName,
					Description: "commit scope lexeme",
					Variants: []*nlpapi.Variant{
						{
							Matcher: &nlpapi.Matcher{Regex: regexp.MustCompile(`\((\w+)\)`)},
						},
					},
					Normalisers: []*nlpapi.NormalisationStep{
						{
							Matcher:    &nlpapi.Matcher{Regex: regexp.MustCompile(`(\w+)`)},
							ReplaceAll: true,
							Formatter: &nlpapi.Formatter{
								Template: "{{ lower . }}",
							},
						},
					},
				},
				{
					Name:        api.CommitBreakingFlagTokenName,
					TokenName:   api.CommitBreakingFlagTokenName,
					Description: "commit scope lexeme",
					Variants: []*nlpapi.Variant{
						{
							Matcher: &nlpapi.Matcher{Regex: regexp.MustCompile(`(!{0,1})`)},
						},
					},
					Normalisers: []*nlpapi.NormalisationStep{
						{
							Matcher:    &nlpapi.Matcher{Regex: regexp.MustCompile(`(!)`)},
							ReplaceAll: true,
						},
					},
				},
			},
		},
		{
			LanguageCode: nlpapi.LanguageCodeWildcard,
			Name:         "issue_tracker_reference",
			Description:  "lexeme to identify issue tracker references from different providers",
			TokenName:    "issue_tracker_reference",
			Variants: []*nlpapi.Variant{
				{
					Name:    "JIRA like issue reference",
					Matcher: &nlpapi.Matcher{Regex: regexp.MustCompile(`[\(\[]{0,1}([\w]{0,6})[-_]([\d]+)[\)\]]{0,1}`)},
					Normalisers: []*nlpapi.NormalisationStep{
						{
							Matcher:    &nlpapi.Matcher{Regex: regexp.MustCompile(`([\w]{0,6})[-_]([\d]+)`)},
							ReplaceAll: true,
							Formatter: &nlpapi.Formatter{
								Template: "{{ upper . }}",
							},
						},
						{
							Matcher:     &nlpapi.Matcher{Regex: regexp.MustCompile(`[-_]`)},
							Replacement: "-",
						},
					},
				},
				{
					Name:    "GitHub issue reference",
					Matcher: &nlpapi.Matcher{Regex: regexp.MustCompile(`[\(\[]{0,1}(#|gh-|GH-)([\d]+)[\)\]]{0,1}`)},
					Normalisers: []*nlpapi.NormalisationStep{
						{
							Name:        "github issue reference",
							Matcher:     &nlpapi.Matcher{Regex: regexp.MustCompile("(#)")},
							Replacement: "gh-",
						},
					},
				},
			},
		},
	}
}

func commitDictionaries() []*nlpapi.Dictionary {
	return []*nlpapi.Dictionary{
		{
			LanguageCode: "en",
			Name:         "weak-feature-signals",
			Description:  "a collection of words that can imply a new feature",
			Weight:       1,
			TokenName:    api.CommitTypeTokenName,
			TokenValue:   typeFeature,
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
			Weight:       2,
			TokenName:    api.CommitTypeTokenName,
			TokenValue:   TypeIgnore,
			Entries: []string{
				"typo",
				"wip",
			},
		},
		{
			LanguageCode: "en",
			Name:         "refactor-signals",
			Weight:       2,
			TokenName:    api.CommitTypeTokenName,
			TokenValue:   typeRefactor,
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
			Name:         "fix-signals",
			Weight:       2,
			TokenName:    api.CommitTypeTokenName,
			TokenValue:   typeFix,
			Entries: []string{
				"fix",
				"bugfix",
				"bug",
			},
		},
		{
			LanguageCode: "en",
			Name:         "docs-signals",
			Weight:       2,
			TokenName:    api.CommitTypeTokenName,
			TokenValue:   typeDocs,
			Entries: []string{
				"document",
				"doc",
			},
		},
		{
			LanguageCode: "en",
			Name:         "test-signals",
			Weight:       2,
			TokenName:    api.CommitTypeTokenName,
			TokenValue:   typeTest,
			Entries: []string{
				"test",
			},
		},
	}
}
