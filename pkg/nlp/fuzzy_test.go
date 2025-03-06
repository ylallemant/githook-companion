package nlp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cases = []struct {
	name  string
	word  string
	entry string
}{
	{
		name:  "equal length",
		word:  "equal",
		entry: "equal",
	},
	{
		name:  "word longer",
		word:  "verylongword",
		entry: "entry",
	},
	{
		name:  "entry longer",
		word:  "word",
		entry: "verylongentry",
	},
	{
		name:  "within entry",
		word:  "thing",
		entry: "somethingsinthemiddle",
	},
	{
		name:  "within spaced entry",
		word:  "thing",
		entry: "some things in the middle",
	},
	{
		name:  "inverted within spaced entry",
		word:  "some things in the middle",
		entry: "thing",
	},
}

func Test_absoluteDistance(t *testing.T) {
	expectations := map[string]float64{
		"equal length":                 0,
		"word longer":                  7,
		"entry longer":                 9,
		"within entry":                 16,
		"within spaced entry":          20,
		"inverted within spaced entry": 20,
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			result := absoluteDistance(c.word, c.entry)

			expected, ok := expectations[c.name]
			assert.Equal(tt, true, ok, fmt.Sprintf("missing test result \"%s\"", c.name))

			if ok {
				assert.Equal(tt, expected, result, "wrong result")
			}
		})
	}
}

func Test_distance(t *testing.T) {
	expectations := map[string]float64{
		"equal length":                 0,
		"word longer":                  9,
		"entry longer":                 11,
		"within entry":                 16,
		"within spaced entry":          20,
		"inverted within spaced entry": 20,
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			result := distance(c.word, c.entry)

			expected, ok := expectations[c.name]
			assert.Equal(tt, true, ok, fmt.Sprintf("missing test result \"%s\"", c.name))

			if ok {
				assert.Equal(tt, expected, result, "wrong result")
			}
		})
	}
}

func Test_basicDistance(t *testing.T) {
	expectations := map[string]float64{
		"equal length":                 5,
		"word longer":                  3,
		"entry longer":                 2,
		"within entry":                 5,
		"within spaced entry":          5,
		"inverted within spaced entry": 5,
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			result := basicDistance(c.word, c.entry)

			expected, ok := expectations[c.name]
			assert.Equal(tt, true, ok, fmt.Sprintf("missing test result \"%s\"", c.name))

			if ok {
				assert.Equal(tt, expected, result, "wrong result")
			}
		})
	}
}

func Test_maxDistance(t *testing.T) {
	expectations := map[string]float64{
		"equal length":                 5,
		"word longer":                  12,
		"entry longer":                 13,
		"within entry":                 21,
		"within spaced entry":          25,
		"inverted within spaced entry": 25,
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			result := maxDistance(c.word, c.entry)

			expected, ok := expectations[c.name]
			assert.Equal(tt, true, ok, fmt.Sprintf("missing test result \"%s\"", c.name))

			if ok {
				assert.Equal(tt, expected, result, "wrong result")
			}
		})
	}
}
