package nlp

import (
	"math"
	"strings"

	edlib "github.com/hbollon/go-edlib"
)

func calculateConfidence(word, entry string) float64 {
	if word == entry {
		return 1.0
	}

	distance := edlib.DamerauLevenshteinDistance(word, entry)

	return float64(distance)
}

func distance(word, entry string) float64 {
	return float64(edlib.DamerauLevenshteinDistance(word, entry))
}

func absoluteDistance(word, entry string) float64 {
	return math.Abs(float64(len(word) - len(entry)))
}

// not sure of this
func basicDistance(word, entry string) float64 {
	return maxDistance(word, entry) - distance(word, entry)
}

func maxDistance(word, entry string) float64 {
	if len(word) < len(entry) {
		return float64(len(entry))
	} else {
		return float64(len(word))
	}
}

// the smaller the better
func indexedDistance(word, entry string) float64 {
	if len(word) > len(entry) {
		return -1.0
	}

	index := strings.Index(entry, word)

	if index > -1 {
		return float64(index)
	}

	return 0.0
}
