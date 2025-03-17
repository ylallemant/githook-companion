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
		entry: "scattered",
	},
	{
		name:  "scattered - sctered",
		word:  "sctered",
		entry: "scattered",
	},
	{
		name:  "scattered - cater",
		word:  "cater",
		entry: "scattered",
	},
	{
		name:  "scattered - ctter",
		word:  "ctter",
		entry: "scattered",
	},
}

// func Test_absoluteDistance(t *testing.T) {
// 	expectations := map[string]float64{
// 		"input longer than entry": 7,
// 		"equal":                   0,
// 		"short - shor":            1,
// 		"short - hort":            1,
// 		"short - sho":             2,
// 		"short - ort":             2,
// 		"short - hor":             2,
// 		"something - somethin":    1,
// 		"something - omething":    1,
// 		"something - somethi":     2,
// 		"something - mething":     2,
// 		"something - ometh":       4,
// 		"something - methi":       4,
// 		"something - ethin":       4,
// 		"scattered - scatterd":    1,
// 		"scattered - sattered":    1,
// 		"scattered - scatted":     2,
// 		"scattered - sctered":     2,
// 		"scattered - cater":       4,
// 		"scattered - ctter":       4,
// 	}

// 	for _, c := range cases {
// 		t.Run(c.name, func(tt *testing.T) {
// 			result := absoluteDistance(c.word, c.entry)

// 			expected, ok := expectations[c.name]
// 			assert.Equal(tt, true, ok, fmt.Sprintf("missing test result \"%s\"", c.name))

// 			if ok {
// 				assert.Equal(tt, expected, result, "wrong result")
// 			}
// 		})
// 	}
// }

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
		"scattered - scatted":     2,
		"scattered - sctered":     2,
		"scattered - cater":       4,
		"scattered - ctter":       4,
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

func Test_maxDistance(t *testing.T) {
	expectations := map[string]float64{
		"input longer than entry": 12,
		"equal":                   5,
		"short - shor":            5,
		"short - hort":            5,
		"short - sho":             5,
		"short - ort":             5,
		"short - hor":             5,
		"something - somethin":    9,
		"something - omething":    9,
		"something - somethi":     9,
		"something - mething":     9,
		"something - ometh":       9,
		"something - methi":       9,
		"something - ethin":       9,
		"scattered - scatterd":    9,
		"scattered - sattered":    9,
		"scattered - scatted":     9,
		"scattered - sctered":     9,
		"scattered - cater":       9,
		"scattered - ctter":       9,
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
		"input longer than entry": 12,
		"equal":                   0,
		"short - shor":            0,
		"short - hort":            1,
		"short - sho":             0,
		"short - ort":             2,
		"short - hor":             1,
		"something - somethin":    0,
		"something - omething":    1,
		"something - somethi":     0,
		"something - mething":     2,
		"something - ometh":       1,
		"something - methi":       2,
		"something - ethin":       3,
		"scattered - scatterd":    9,
		"scattered - sattered":    9,
		"scattered - scatted":     9,
		"scattered - sctered":     9,
		"scattered - cater":       9,
		"scattered - ctter":       9,
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			result := positionDistance(c.word, c.entry)

			expected, ok := expectations[c.name]
			assert.Equal(tt, true, ok, fmt.Sprintf("missing test result \"%s\"", c.name))

			if ok {
				assert.Equal(tt, expected, float64(result), "wrong result")
			}
		})
	}
}

func Test_basicDistance(t *testing.T) {
	expectations := map[string]float64{
		"input longer than entry": 0.875,
		"equal":                   0,
		"short - shor":            0.1,
		"short - hort":            0.2,
		"short - sho":             0.2,
		"short - ort":             0.4,
		"short - hor":             0.3,
		"something - somethin":    0.05555555555555555,
		"something - omething":    0.1111111111111111,
		"something - somethi":     0.1111111111111111,
		"something - mething":     0.2222222222222222,
		"something - ometh":       0.2777777777777778,
		"something - methi":       0.3333333333333333,
		"something - ethin":       0.3888888888888889,
		"scattered - scatterd":    0.5555555555555556,
		"scattered - sattered":    0.5555555555555556,
		"scattered - scatted":     0.6111111111111112,
		"scattered - sctered":     0.6111111111111112,
		"scattered - cater":       0.7222222222222222,
		"scattered - ctter":       0.7222222222222222,
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

// func Test_scatteredDistance(t *testing.T) {
// 	expectations := map[string]float64{
// 		"input longer than entry": 1.2352941176470589,
// 		"equal":                   0,
// 		"short - shor":            0,
// 		"short - hort":            0,
// 		"short - sho":             0,
// 		"short - ort":             0.2222222222222222,
// 		"short - hor":             0,
// 		"something - somethin":    0.3,
// 		"something - omething":    0.2,
// 		"something - somethi":     0.2,
// 		"something - mething":     0.2,
// 		"something - ometh":       0,
// 		"something - methi":       0.1,
// 		"something - ethin":       0.12,
// 		"scattered - scatterd":    0.6666666666666666,
// 		"scattered - sattered":    0.6666666666666666,
// 		"scattered - scatted":     0.5,
// 		"scattered - sctered":     0.5384615384615384,
// 		"scattered - cater":       0.2727272727272727,
// 		"scattered - ctter":       0.17647058823529413,
// 	}

// 	for _, c := range cases {
// 		t.Run(c.name, func(tt *testing.T) {
// 			result := scatteredDistance(c.word, c.entry)

// 			expected, ok := expectations[c.name]
// 			assert.Equal(tt, true, ok, fmt.Sprintf("missing test result \"%s\"", c.name))

// 			if ok {
// 				assert.Equal(tt, expected, result, "wrong result")
// 			}
// 		})
// 	}
// }

func Test_calculateConfidence(t *testing.T) {
	expectations := map[string]float64{
		"input longer than entry": 0,
		"equal":                   1,
		"short - shor":            0.9,
		"short - hort":            0.8,
		"short - sho":             0.8,
		"short - ort":             0.6,
		"short - hor":             0.7,
		"something - somethin":    0.9444444444444444,
		"something - omething":    0.8888888888888888,
		"something - somethi":     0.8888888888888888,
		"something - mething":     0.7777777777777778,
		"something - ometh":       0.7222222222222222,
		"something - methi":       0.6666666666666667,
		"something - ethin":       0.6111111111111112,
		"scattered - scatterd":    0.4444444444444444,
		"scattered - sattered":    0.4444444444444444,
		"scattered - scatted":     0.38888888888888884,
		"scattered - sctered":     0.38888888888888884,
		"scattered - cater":       0.2777777777777778,
		"scattered - ctter":       0.2777777777777778,
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

func Test_customTest(t *testing.T) {
	cases := []struct {
		word     string
		entry    string
		expected float64
	}{
		{entry: "teller", word: "tes", expected: 0.8333333333333334},
		{entry: "test", word: "es", expected: 0.375},
		{entry: "testing", word: "es", expected: 0.42857142857142855},
		{entry: "我愛你", word: "我叫你", expected: 0.5555555555555556},
		{entry: "boxer", word: "rexob", expected: 0.9},
		{entry: "boxer", word: "reoxb", expected: 0.9},
		{entry: "boxer", word: "robex", expected: 0.8},
		{entry: "butterfly", word: "cow", expected: 1},
		{entry: "butterfly", word: "butterfly", expected: 0},
		{entry: "butterfly", word: "butter", expected: 0.16666666666666666},
		{entry: "butterfly", word: "fly", expected: 0.6666666666666666},
		{entry: "butterfly", word: "terfly", expected: 0.3333333333333333},
		{entry: "werl", word: "berl", expected: 0.625},
		{entry: "berlin", word: "berl", expected: 0.16666666666666666},
		{entry: "Bürger", word: "burger", expected: 0.6428571428571429},
		{entry: "fly", word: "butterfly", expected: 0.8333333333333334},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("search %s in %s", c.word, c.entry), func(tt *testing.T) {
			result := basicDistance(c.word, c.entry)

			assert.Equal(tt, c.expected, result, "wrong result")
		})
	}
}
