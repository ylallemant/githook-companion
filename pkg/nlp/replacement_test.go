package nlp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_searchStringToRegex(t *testing.T) {
	cases := []struct {
		name     string
		search   string
		expected string
	}{
		{
			name:     "empty search",
			search:   "",
			expected: "",
		},
		{
			name:     "no special characters",
			search:   "word",
			expected: "\\b(word)\\b",
		},
		{
			name:     "special character in the middle",
			search:   "one-two",
			expected: "\\b(one-two)\\b",
		},
		{
			name:     "special character left",
			search:   "(one-two",
			expected: "(\\(one-two)\\b",
		},
		{
			name:     "special character right",
			search:   "one-two)",
			expected: "\\b(one-two\\))",
		},
		{
			name:     "between special characters",
			search:   "(one-two)",
			expected: "(\\(one-two\\))",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			regex := searchStringToRegex(c.search)

			if regex != nil {
				assert.Equal(tt, c.expected, regex.String(), "wrong regex")
			}
		})
	}
}

func Test_secureReplaceAllString(t *testing.T) {
	cases := []struct {
		name     string
		text     string
		old      string
		new      string
		expected string
	}{
		{
			name:     "empty search",
			text:     "not important",
			old:      "word",
			new:      "bar",
			expected: "not important",
		},
		{
			name:     "empty text",
			text:     "",
			old:      "word",
			new:      "bar",
			expected: "",
		},
		{
			name:     "no occurence",
			text:     "nothing to find",
			old:      "word",
			new:      "bar",
			expected: "nothing to find",
		},
		{
			name:     "at the start",
			text:     "word to please",
			old:      "word",
			new:      "bar",
			expected: "bar to please",
		},
		{
			name:     "at the end",
			text:     "final word",
			old:      "word",
			new:      "bar",
			expected: "final bar",
		},
		{
			name:     "single occurence",
			text:     "single word",
			old:      "word",
			new:      "bar",
			expected: "single bar",
		},
		{
			name:     "two occurences",
			text:     "single word but with no mixword, only one word more",
			old:      "word",
			new:      "bar",
			expected: "single bar but with no mixword, only one bar more",
		},
		{
			name:     "multiple occurences",
			text:     "a word is a word, is a word, given by a wordless word in the world",
			old:      "word",
			new:      "bar",
			expected: "a bar is a bar, is a bar, given by a wordless bar in the world",
		},
		{
			name:     "multiple occurences with longer replacement",
			text:     "a word is a word, is a word, given by a wordless word in the world",
			old:      "word",
			new:      "barbarian",
			expected: "a barbarian is a barbarian, is a barbarian, given by a wordless barbarian in the world",
		},
		{
			name:     "special characters",
			text:     "single (word) but with no mixword, only one (word) more",
			old:      "(word)",
			new:      "bar",
			expected: "single bar but with no mixword, only one bar more",
		},
		{
			name:     "special characters at start",
			text:     "[word] to please",
			old:      "[word]",
			new:      "bar",
			expected: "bar to please",
		},
		{
			name:     "special characters at end",
			text:     "final (word)",
			old:      "(word)",
			new:      "bar",
			expected: "final bar",
		},
		{
			name:     "with spaces",
			text:     "single word in the middle",
			old:      "word",
			new:      "bar",
			expected: "single bar in the middle",
		},
		{
			name:     "ignore in string occurence",
			text:     "mixword",
			old:      "word",
			new:      "bar",
			expected: "mixword",
		},
		{
			name:     "ignore any in-string occurence",
			text:     "wordwordwordword",
			old:      "word",
			new:      "bar",
			expected: "wordwordwordword",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			valid := secureReplaceAllString(c.text, c.old, c.new)

			assert.Equal(tt, c.expected, valid, "wrong result")
		})
	}
}
