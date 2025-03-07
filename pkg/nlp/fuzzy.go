package nlp

import (
	"fmt"
	"math"
	"strings"

	edlib "github.com/hbollon/go-edlib"
)

func calculateConfidence(word, entry string) float64 {

	if word == entry {
		return 1.0
	}

	if len(word) > len(entry) {
		return 0
	}

	position := positionDistance(word, entry)
	characters := 1.1

	if position > -1 {
		characters = basicDistance(word, entry)
	} else {
		// no clear match
		characters = charDistance(word, entry)
	}

	if characters > 1 {
		return 0
	}

	return float64(1 - characters)
}

func distance(word, entry string) float64 {
	return float64(edlib.DamerauLevenshteinDistance(word, entry))
}

func absoluteDistance(word, entry string) float64 {
	return math.Abs(float64(len(word) - len(entry)))
}

// not sure of this
func basicDistance(word, entry string) float64 {
	position := positionDistance(word, entry)
	maximum := maxDistance(word, entry)
	absolute := distance(word, entry)

	if position < 0 {
		return 1
	}

	lengthRatio := float64(absolute) / float64(maximum)
	positionRatio := absolute

	if absolute > 0 {
		positionRatio = float64(position) / float64(absolute)
	}

	fmt.Println("  - position:           ", position)
	fmt.Println("  - maximum:            ", maximum)
	fmt.Println("  - absolute:           ", absolute)
	fmt.Println("  - position / maxPos:  ", positionRatio)
	fmt.Println("  - absolute / maximum:  ", float64(absolute)/float64(maximum))
	fmt.Println("  - result:              ", lengthRatio-positionRatio)

	return math.Abs(lengthRatio - positionRatio)
}

func maxDistance(word, entry string) float64 {
	if len(word) < len(entry) {
		return float64(len(entry))
	} else {
		return float64(len(word))
	}
}

// positionDistance checks the position of the word in the entry
// 0 is best, if not found -1 is returned
func positionDistance(word, entry string) float64 {
	index := strings.Index(entry, word)

	return float64(index)
}

// charDistance checks for the position of the individual characters
// of the word in the entry.
// the lower the returned value, the better
func charDistance(word, entry string) float64 {
	charPositions := make([][]int, 0)
	maxPosition := len(entry) - 1

	for _, rune := range word {
		char := string(rune)
		position := 0
		lastPosition := 0
		positions := make([]int, 0)

		for position > -1 && position < maxPosition {

			// position in truncated entry
			position = strings.Index(entry[position:], char)

			if position < 0 {
				// char not found, stop search
				break
			}

			// calcule real position in entry
			position = lastPosition + position

			positions = append(positions, position)

			// calculate next start position
			position = position + 1
			lastPosition = position
		}

		charPositions = append(charPositions, positions)
	}

	// make a sums per index for:
	//    - distance of scatered characters
	//    - sum of all positions
	//    - missing character count
	//
	// characters following each other in the right order
	// do not increase the distance (perfect match)
	//
	// ex character positions :
	// [[5 15] [6 16] [7 12 20] [8 13] [9]]
	// ex sum by index:
	// [
	//   map[distance:5 missed:0 sum:35],
	//   map[distance:27 missed:0 sum:56]
	// ]
	distances := make([]map[string]int, 0)

	index := 0
	maxLength := 1

	for index < maxLength {
		lastPosition := 0
		totalDinstance := 0
		totalSum := 0
		totalMissed := 0
		information := make(map[string]int)

		for _, positions := range charPositions {
			if index < len(positions) {
				totalSum = totalSum + positions[index]

				if positions[index] == lastPosition+1 {
					lastPosition = positions[index]
					continue
				}

				lastPosition = positions[index]
				totalDinstance = totalDinstance + lastPosition
			} else if index == 0 {
				totalMissed = totalMissed + 1
			}

			if index == 0 && maxLength < len(positions) {
				maxLength = len(positions)
			}
		}

		information["distance"] = totalDinstance
		information["missed"] = totalMissed
		information["sum"] = totalSum

		distances = append(distances, information)
		index = index + 1
	}

	if len(distances) == 0 && distances[0]["distance"] == 0 {
		// strings are equal
		// return no distance as 0
		return 0
	}

	// get a sum of the distances and sums
	// example:
	// [
	//   map[distance:5 missed:0 sum:35],
	//   map[distance:27 missed:0 sum:56]
	// ]
	//
	// sum distance 32
	// sum sum      91
	// 32 / 91 = 0.3516483516483517
	distance := 0
	sum := 0

	for _, information := range distances {
		distance = distance + information["distance"] + information["missed"]
		sum = sum + information["sum"]
	}

	// return distance / sum ratio
	return float64(distance) / float64(sum)
}
