package validate

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	cases := []struct {
		name         string
		message      string
		expected     string
		expectError  bool
		errorMessage string
	}{
		{
			name:        "not able to detect commit type",
			message:     "some changes",
			expected:    "",
			expectError: true,
			errorMessage: `commit message malformed
  you didn't commit on the command line, commit type can not be added interactively
  please make sure to provide a commit type prefix in your message
  format: "<commit-type-prefix>: <commit-message>"
  available commit types:
    - feat: a new feature is introduced with the changes
    - fix: a bug fix has occurred
    - docs: updates to documentation such as a the README or other markdown files
    - test: including new or correcting previous tests
    - refactor: refactored code that neither fixes a bug nor adds a feature

			`,
		},
		{
			name:        "detect commit type from dictionary with plain value",
			message:     "add some changes",
			expected:    "feat: add some changes\n",
			expectError: false,
		},
		{
			name:        "detect commit type from dictionary with synonym",
			message:     "added some changes",
			expected:    "feat: add some changes\n",
			expectError: false,
		},
		{
			name:        "detect commit type from existsing type",
			message:     "refactor: some changes",
			expected:    "refactor: some changes\n",
			expectError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			b := bytes.NewBufferString("")

			rootCmd.SetOut(b)
			rootCmd.SetArgs([]string{"-m", c.message})
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
