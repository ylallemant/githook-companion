package commit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ylallemant/githooks-butler/pkg/api"
)

func TestMessageType(t *testing.T) {
	cases := []struct {
		name               string
		message            string
		config             *api.Config
		expectedCommitType string
		expectedTypeFound  bool
	}{
		{
			name:    "no matching commit type prefix",
			message: "some changes",
			config: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{
						{
							Type: "feat",
						},
					},
				},
			},
			expectedCommitType: "",
			expectedTypeFound:  false,
		},
		{
			name:    "matching commit type always lower case",
			message: "feat:some changes",
			config: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{
						{
							Type: "FEAT",
						},
					},
				},
			},
			expectedCommitType: "feat",
			expectedTypeFound:  true,
		},
		{
			name:    "matching commit type from multiple",
			message: "Docs:some changes",
			config: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{
						{
							Type: "feat",
						},
						{
							Type: "DOCS",
						},
						{
							Type: "test",
						},
					},
				},
			},
			expectedCommitType: "docs",
			expectedTypeFound:  true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			commitType, found := messageType(c.message, c.config)

			assert.Equal(tt, c.expectedCommitType, commitType, "wrong result")
			assert.Equal(tt, c.expectedTypeFound, found, "wrong result")
		})
	}
}

func TestEnsureFormat(t *testing.T) {
	cases := []struct {
		name       string
		message    string
		commitType string
		expected   string
	}{
		{
			name:       "add commit type prefix",
			message:    "some changes",
			commitType: "feat",
			expected:   "feat: some changes",
		},
		{
			name:       "dont change commit message case",
			message:    "Some Heads-UP changes",
			commitType: "feat",
			expected:   "feat: Some Heads-UP changes",
		},
		{
			name:       "ensure lower case commit type prefix",
			message:    "Some Heads-UP changes",
			commitType: "FEAT",
			expected:   "feat: Some Heads-UP changes",
		},
		{
			name:       "ensure lower case commit type prefix",
			message:    "Some Heads-UP changes",
			commitType: "FEAT",
			expected:   "feat: Some Heads-UP changes",
		},
		{
			name:       "ensure lower case commit type prefix on existing prefix",
			message:    "FEAT : Some Heads-UP changes",
			commitType: "FEAT",
			expected:   "feat: Some Heads-UP changes",
		},
		{
			name:       "ignore missing colon on existing prefix",
			message:    "FEAT  Some Heads-UP FEAT changes",
			commitType: "FEAT",
			expected:   "feat: Some Heads-UP FEAT changes",
		},
		{
			name:       "only change first commit type name occurance",
			message:    "FEAT : Some Heads-UP FEAT changes",
			commitType: "FEAT",
			expected:   "feat: Some Heads-UP FEAT changes",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			result := EnsureFormat(c.message, c.commitType)

			assert.Equal(tt, c.expected, result, "wrong result")
		})
	}
}

func TestEnsureDictionaryValue(t *testing.T) {
	cases := []struct {
		name       string
		message    string
		dictionary *api.CommitTypeDictionary
		expected   string
	}{
		{
			name:    "replace first word",
			message: "added some changes",
			dictionary: &api.CommitTypeDictionary{
				Name:  "add",
				Value: "add",
				Type:  "feat",
				Synonyms: []string{
					"adds",
					"added",
					"adding",
					"new",
				},
			},
			expected: "add some changes",
		},
		{
			name:    "ignore upfront spaces",
			message: "    added some changes",
			dictionary: &api.CommitTypeDictionary{
				Name:  "add",
				Value: "add",
				Type:  "feat",
				Synonyms: []string{
					"adds",
					"added",
					"adding",
					"new",
				},
			},
			expected: "add some changes",
		},
		{
			name:    "ignore case",
			message: "AdDs some changes",
			dictionary: &api.CommitTypeDictionary{
				Name:  "add",
				Value: "add",
				Type:  "feat",
				Synonyms: []string{
					"adds",
					"added",
					"adding",
					"new",
				},
			},
			expected: "add some changes",
		},
		{
			name:    "ignore punctuation",
			message: "AdDs: some changes",
			dictionary: &api.CommitTypeDictionary{
				Name:  "add",
				Value: "add",
				Type:  "feat",
				Synonyms: []string{
					"adds",
					"added",
					"adding",
					"new",
				},
			},
			expected: "add: some changes",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			result := EnsureDictionaryValue(c.message, c.dictionary)

			assert.Equal(tt, c.expected, result, "wrong result")
		})
	}
}
