package nlp

import (
	"fmt"
	"os"
	"testing"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
	"github.com/stretchr/testify/assert"
)

func TestTokernizer(t *testing.T) {
	cases := []struct {
		name         string
		message      string
		expected     string
		expectError  bool
		errorMessage string
	}{
		{
			name:        "stemming working",
			message:     "working",
			expected:    "work",
			expectError: false,
		},
		{
			name:        "stemming worked",
			message:     "worked",
			expected:    "work",
			expectError: false,
		},
		{
			name:        "stemming works",
			message:     "works",
			expected:    "work",
			expectError: false,
		},
		{
			name:        "stemming adding",
			message:     "adding",
			expected:    "add",
			expectError: false,
		},
		{
			name:        "stemming added",
			message:     "AdDed",
			expected:    "add",
			expectError: false,
		},
		{
			name:        "stemming adds",
			message:     "adds",
			expected:    "add",
			expectError: false,
		},
	}

	res, _ := en.New().GetResource()
	file, _ := os.OpenFile("en.txt", os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()

	file.Write(res)

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			lemmatizer, err := golem.New(en.New())
			assert.Nil(tt, err)

			assert.Nil(tt, err)

			word := lemmatizer.Lemma(c.message)
			fmt.Println(lemmatizer.Lemmas(c.message))

			if c.expectError {
				assert.NotNil(tt, err)
				assert.Equal(tt, c.errorMessage, err.Error(), "wrong error massage")
			} else {
				assert.Nil(tt, err)
				assert.Equal(tt, c.expected, word, "wrong result")
			}
		})
	}
}
