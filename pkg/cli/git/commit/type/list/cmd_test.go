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
			expected: `┌────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│ Commit Types                                                                                               │
├──────────┬───────────────────────────────────────────────────────────────────────────────────┬─────────────┤
│ TYPE     │ DESCRIPTION                                                                       │ AUTO-FORMAT │
├──────────┼───────────────────────────────────────────────────────────────────────────────────┼─────────────┤
│ feat     │ a new feature is introduced with the changes                                      │ true        │
│ refactor │ refactored code that neither fixes a bug nor adds a feature                       │ true        │
│ ignore   │ commit can be ignored by other tools                                              │ false       │
│ fix      │ a bug fix has been implemented                                                    │ true        │
│ docs     │ documentation only changes                                                        │ true        │
│ test     │ including new or correcting previous tests                                        │ true        │
│ perf     │ a code change that improves performance                                           │ true        │
│ style    │ changes that do not affect the meaning of the code (white-space, formatting, ...) │ true        │
│ chore    │ other changes that don't modify src or test files                                 │ true        │
│ build    │ changes that affect the build system or external dependencies                     │ true        │
│ ci       │ changes to CI configuration files and scripts                                     │ true        │
└──────────┴───────────────────────────────────────────────────────────────────────────────────┴─────────────┘
`,
			expectError: false,
		},
		// disabled along with the "-c" flag
		// 		{
		// 			name:       "with simple config",
		// 			configPath: "../../../../../../test/configs/simple/config.yaml",
		// 			expected: `┌───────────────────────────────────────────────────────────────────┐
		// │ Commit Types                                                      │
		// ├──────┬──────────────────────────────────────────────┬─────────────┤
		// │ TYPE │ DESCRIPTION                                  │ AUTO-FORMAT │
		// ├──────┼──────────────────────────────────────────────┼─────────────┤
		// │ feat │ a new feature is introduced with the changes │ true        │
		// │ fix  │ a bug fix has occurred                       │ true        │
		// └──────┴──────────────────────────────────────────────┴─────────────┘
		// `,
		// 			expectError: false,
		// 		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			b := bytes.NewBufferString("")
			args := []string{}
			args = append(args, "--fallback-config")

			rootCmd.SetOut(b)
			if c.configPath != "" {
				args = append(args, []string{"-c", c.configPath}...)
			}

			rootCmd.SetArgs(args)

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
