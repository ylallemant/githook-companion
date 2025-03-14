package commit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ylallemant/githook-companion/pkg/config"
)

func TestEnsureFormat(t *testing.T) {
	cases := []struct {
		name     string
		message  string
		expected string
	}{
		// {
		// 	name:       "add commit type prefix",
		// 	message:    "some changes",
		// 	commitType: "feat",
		// 	expected:   "feat: some changes",
		// },
		{
			name:     "ensure lower case commit message",
			message:  "Some Heads-UP added",
			expected: "FEAT: some heads-up added",
		},
		{
			name:     "ensure upper case commit type prefix",
			message:  "Some Heads-UP changes",
			expected: "REFACTOR: some heads-up changes",
		},
		{
			name:     "ensure upper case commit type prefix on existing prefix",
			message:  "feat : Some Heads-UP changes",
			expected: "FEAT: some heads-up changes",
		},
		{
			name:     "ignore missing colon on existing prefix",
			message:  "FEAT  Some Heads-UP changes",
			expected: "FEAT: some heads-up changes",
		},
		{
			name:     "only change first commit type name occurance",
			message:  "feAt : Some Heads-UP FEAT changes",
			expected: "FEAT: some heads-up feat changes",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			cfg := config.Default()
			_, _, commitTypeToken, tokens := Validate(c.message, cfg)

			result, err := EnsureFormat(c.message, cfg.Commit.MessageTemplate, commitTypeToken, tokens)
			assert.Nil(tt, err)

			assert.Equal(tt, c.expected, result, "wrong result")
		})
	}
}
