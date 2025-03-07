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
		name:  "input longer than entry",
		word:  "verylongword",
		entry: "entry",
	},
	{
		name:  "equal",
		word:  "equal",
		entry: "equal",
	},
	{
		name:  "short - shor",
		word:  "shor",
		entry: "short",
	},
	{
		name:  "short - hort",
		word:  "hort",
		entry: "short",
	},
	{
		name:  "short - sho",
		word:  "sho",
		entry: "short",
	},
	{
		name:  "short - ort",
		word:  "ort",
		entry: "short",
	},
	{
		name:  "short - hor",
		word:  "hor",
		entry: "short",
	},
	{
		name:  "something - somethin",
		word:  "somethin",
		entry: "something",
	},
	{
		name:  "something - omething",
		word:  "omething",
		entry: "something",
	},
	{
		name:  "something - somethi",
		word:  "somethi",
		entry: "something",
	},
	{
		name:  "something - mething",
		word:  "mething",
		entry: "something",
	},
	{
		name:  "something - ometh",
		word:  "ometh",
		entry: "something",
	},
	{
		name:  "something - methi",
		word:  "methi",
		entry: "something",
	},
	{
		name:  "something - ethin",
		word:  "ethin",
		entry: "something",
	},
	{
		name:  "scattered - scatterd",
		word:  "scatterd",
		entry: "scattered",
	},
	{
		name:  "scattered - sattered",
		word:  "sattered",
		entry: "scattered",
	},
	{
		name:  "scattered - scatted",
		word:  "scatted",
		entry: "something",
	},
	{
		name:  "scattered - sctered",
		word:  "sctered",
		entry: "something",
	},
	{
		name:  "scattered - cater",
		word:  "cater",
		entry: "something",
	},
	{
		name:  "scattered - ctter",
		word:  "ctter",
		entry: "something",
	},
}

func Test_absoluteDistance(t *testing.T) {
	expectations := map[string]float64{
		"input longer than entry": 7,
		"equal":                   0,
		"short - shor":            1,
		"short - hort":            1,
		"short - sho":             2,
		"short - ort":             2,
		"short - hor":             2,
		"something - somethin":    1,
		"something - omething":    1,
		"something - somethi":     2,
		"something - mething":     2,
		"something - ometh":       4,
		"something - methi":       4,
		"something - ethin":       4,
		"scattered - scatterd":    1,
		"scattered - sattered":    1,
		"scattered - scatted":     2,
		"scattered - sctered":     2,
		"scattered - cater":       4,
		"scattered - ctter":       4,
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
		"input longer than entry": 9,
		"equal":                   0,
		"short - shor":            1,
		"short - hort":            1,
		"short - sho":             2,
		"short - ort":             2,
		"short - hor":             2,
		"something - somethin":    1,
		"something - omething":    1,
		"something - somethi":     2,
		"something - mething":     2,
		"something - ometh":       4,
		"something - methi":       4,
		"something - ethin":       4,
		"scattered - scatterd":    1,
		"scattered - sattered":    1,
		"scattered - scatted":     7,
		"scattered - sctered":     7,
		"scattered - cater":       8,
		"scattered - ctter":       8,
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
		"input longer than entry": 1,
		"equal":                   0,
		"short - shor":            0.2,
		"short - hort":            0.8,
		"short - sho":             0.4,
		"short - ort":             0.6,
		"short - hor":             0.6,
		"something - somethin":    9,
		"something - omething":    5,
		"something - somethi":     2,
		"something - mething":     20,
		"something - ometh":       20,
		"something - methi":       20,
		"something - ethin":       20,
		"scattered - scatterd":    1,
		"scattered - sattered":    1,
		"scattered - scatted":     1,
		"scattered - sctered":     1,
		"scattered - cater":       1,
		"scattered - ctter":       1,
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
		"input longer than entry": 0,
		"equal":                   7,
		"short - shor":            9,
		"short - hort":            5,
		"short - sho":             16,
		"short - ort":             20,
		"short - hor":             20,
		"something - somethin":    9,
		"something - omething":    5,
		"something - somethi":     2,
		"something - mething":     20,
		"something - ometh":       20,
		"something - methi":       20,
		"something - ethin":       20,
		"scattered - scatterd":    9,
		"scattered - sattered":    5,
		"scattered - scatted":     16,
		"scattered - sctered":     20,
		"scattered - cater":       20,
		"scattered - ctter":       20,
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

func Test_positionDistance(t *testing.T) {
	expectations := map[string]float64{
		"input longer than entry": 0,
		"equal":                   7,
		"short - shor":            9,
		"short - hort":            5,
		"short - sho":             16,
		"short - ort":             20,
		"short - hor":             20,
		"something - somethin":    9,
		"something - omething":    5,
		"something - somethi":     2,
		"something - mething":     20,
		"something - ometh":       20,
		"something - methi":       20,
		"something - ethin":       20,
		"scattered - scatterd":    9,
		"scattered - sattered":    5,
		"scattered - scatted":     16,
		"scattered - sctered":     20,
		"scattered - cater":       20,
		"scattered - ctter":       20,
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			result := positionDistance(c.word, c.entry)

			expected, ok := expectations[c.name]
			assert.Equal(tt, true, ok, fmt.Sprintf("missing test result \"%s\"", c.name))

			if ok {
				assert.Equal(tt, expected, result, "wrong result")
			}
		})
	}
}

func Test_charDistance(t *testing.T) {
	expectations := map[string]float64{
		"input longer than entry": 0,
		"equal":                   7,
		"short - shor":            9,
		"short - hort":            5,
		"short - sho":             16,
		"short - ort":             20,
		"short - hor":             20,
		"something - somethin":    9,
		"something - omething":    5,
		"something - somethi":     2,
		"something - mething":     20,
		"something - ometh":       20,
		"something - methi":       20,
		"something - ethin":       20,
		"scattered - scatterd":    9,
		"scattered - sattered":    5,
		"scattered - scatted":     16,
		"scattered - sctered":     20,
		"scattered - cater":       20,
		"scattered - ctter":       20,
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			result := charDistance(c.word, c.entry)

			expected, ok := expectations[c.name]
			assert.Equal(tt, true, ok, fmt.Sprintf("missing test result \"%s\"", c.name))

			if ok {
				assert.Equal(tt, expected, result, "wrong result")
			}
		})
	}
}

func Test_calculateConfidence(t *testing.T) {
	expectations := map[string]float64{
		"input longer than entry": 0,
		"equal":                   7,
		"short - shor":            9,
		"short - hort":            5,
		"short - sho":             16,
		"short - ort":             20,
		"short - hor":             20,
		"something - somethin":    9,
		"something - omething":    5,
		"something - somethi":     2,
		"something - mething":     20,
		"something - ometh":       20,
		"something - methi":       20,
		"something - ethin":       20,
		"scattered - scatterd":    9,
		"scattered - sattered":    5,
		"scattered - scatted":     16,
		"scattered - sctered":     20,
		"scattered - cater":       20,
		"scattered - ctter":       20,
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			result := calculateConfidence(c.word, c.entry)

			expected, ok := expectations[c.name]
			assert.Equal(tt, true, ok, fmt.Sprintf("missing test result \"%s\"", c.name))

			if ok {
				assert.Equal(tt, expected, result, "wrong result")
			}
		})
	}
}
