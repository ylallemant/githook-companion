package list

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	cases := []struct {
		name         string
		configPath   string
		expected     string
		expectError  bool
		errorMessage string
	}{
		{
			name:       "with default config",
			configPath: "",
			expected: `┌──────────────────────────────────────────────────────────────────────────────────┐
│ Commit Types                                                                     │
├──────────┬───────────────────────────────────────────────────────────────────────┤
│ TYPE     │ DESCRIPTION                                                           │
├──────────┼───────────────────────────────────────────────────────────────────────┤
│ feat     │ a new feature is introduced with the changes                          │
│ fix      │ a bug fix has occurred                                                │
│ docs     │ updates to documentation such as a the README or other markdown files │
│ test     │ including new or correcting previous tests                            │
│ refactor │ refactored code that neither fixes a bug nor adds a feature           │
│ breaking │ introducing a breaking change in input or output behaviour            │
└──────────┴───────────────────────────────────────────────────────────────────────┘
`,
			expectError: false,
		},
		{
			name:       "with simple config",
			configPath: "../../../../../test/configs/simple.yaml",
			expected: `┌─────────────────────────────────────────────────────┐
│ Commit Types                                        │
├──────┬──────────────────────────────────────────────┤
│ TYPE │ DESCRIPTION                                  │
├──────┼──────────────────────────────────────────────┤
│ feat │ a new feature is introduced with the changes │
│ fix  │ a bug fix has occurred                       │
└──────┴──────────────────────────────────────────────┘
`,
			expectError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			b := bytes.NewBufferString("")

			rootCmd.SetOut(b)
			if c.configPath != "" {
				rootCmd.SetArgs([]string{"-c", c.configPath})
			}

			cmdErr := rootCmd.Execute()

			out, err := io.ReadAll(b)
			assert.Nil(tt, err)

			if c.expectError {
				assert.NotNil(tt, cmdErr)
				assert.Equal(tt, c.errorMessage, cmdErr.Error(), "wrong error massage")
			} else {
				assert.Nil(tt, cmdErr)
				assert.Equal(tt, c.expected, string(out), "wrong result")
			}
		})
	}
}
