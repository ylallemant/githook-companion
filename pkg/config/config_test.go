package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ylallemant/githooks-butler/pkg/api"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name         string
		cfgA         *api.Config
		cfgB         *api.Config
		expected     *api.Config
		throwsError  bool
		errorMessage string
	}{
		{
			name: "merge into empty config",
			cfgA: &api.Config{},
			cfgB: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{
						{
							Type: typeFeature,
						},
					},
				},
			},
			expected: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{
						{
							Type: typeFeature,
						},
					},
					Dictionaries: []*api.CommitTypeDictionary{},
				},
				Dependencies: []*api.Dependency{},
			},
		},
		{
			name: "merge empty into config",
			cfgA: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{
						{
							Type: typeFeature,
						},
					},
				},
			},
			cfgB: &api.Config{},
			expected: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{
						{
							Type: typeFeature,
						},
					},
					Dictionaries: []*api.CommitTypeDictionary{},
				},
				Dependencies: []*api.Dependency{},
			},
		},
		{
			name: "merge non emty configs",
			cfgA: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{
						{
							Type:        typeFeature,
							Description: "a new feature is introduced with the changes",
						},
					},
					DefaultType: "feat",
				},
			},
			cfgB: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{
						{
							Type:        typeDocs,
							Description: "updates to documentation such as a the README or other markdown files",
						},
					},
					DefaultType: "docs",
				},
			},
			expected: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{
						{
							Type:        typeFeature,
							Description: "a new feature is introduced with the changes",
						},
						{
							Type:        typeDocs,
							Description: "updates to documentation such as a the README or other markdown files",
						},
					},
					Dictionaries: []*api.CommitTypeDictionary{},
					DefaultType:  "docs",
				},
				Dependencies: []*api.Dependency{},
			},
		},
		{
			name: "merge with deep overwrites",
			cfgA: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{
						{
							Type:        typeFeature,
							Description: "a new feature is introduced with the changes",
						},
					},
					DefaultType: "feat",
				},
			},
			cfgB: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{
						{
							Type:        typeFeature,
							Description: "overwritten text",
						},
						{
							Type:        typeDocs,
							Description: "updates to documentation such as a the README or other markdown files",
						},
					},
					DefaultType: "docs",
				},
				Dependencies: []*api.Dependency{},
			},
			expected: &api.Config{
				Commit: &api.Commit{
					Types: []*api.CommitType{
						{
							Type:        typeFeature,
							Description: "overwritten text",
						},
						{
							Type:        typeDocs,
							Description: "updates to documentation such as a the README or other markdown files",
						},
					},
					Dictionaries: []*api.CommitTypeDictionary{},
					DefaultType:  "docs",
				},
				Dependencies: []*api.Dependency{},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			merged, err := Merge(c.cfgA, c.cfgB)

			if c.throwsError {
				assert.NotNil(tt, err)
				assert.Equal(tt, c.errorMessage, err.Error(), "wrong error massage")
			} else {
				assert.Nil(tt, err)
			}

			assert.Equal(tt, *c.expected, *merged, "wrong result")
		})
	}
}
