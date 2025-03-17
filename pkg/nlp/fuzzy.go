package nlp

import (
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

	characters := basicDistance(word, entry)

	return float64(1 - characters)
}

func distance(word, entry string) float64 {
	return float64(edlib.DamerauLevenshteinDistance(word, entry))
}

// func absoluteDistance(word, entry string) float64 {
// 	return math.Abs(float64(len(word) - len(entry)))
// }

// maxDistance returns the maximum possible distance
// between word and entry
func maxDistance(word, entry string) float64 {
	if len(word) < len(entry) {
		return float64(len(entry))
	} else {
		return float64(len(word))
	}
}

// positionDistance checks the position of the word in the entry
// 0 is best, if not found -1 is returned
func positionDistance(word, entry string) int {
	index := strings.Index(entry, word)

	if index < 0 {
		return int(maxDistance(word, entry))
	}

	return index
}

/*
basicDistance returns a float between 0 and 1.
The lower the more the compared terms are similar.
On strict equality returns 0.
On complete difference returns 1.
*/
func basicDistance(word, entry string) float64 {
	position := float64(positionDistance(word, entry))
	distance := distance(word, entry)
	maximum := maxDistance(word, entry)

	return float64(position+distance) / float64(maximum*2)
}

// scatteredDistance checks for the position of the individual characters
// of the word in the entry.
// the lower the returned value, the better
// func scatteredDistance(word, entry string) float64 {
// 	maxDistance := int(maxDistance(word, entry))
// 	charPositions := make([][]int, 0)
// 	maxPosition := len(entry) - 1

// 	firstCharPositions := make(map[string]int)

// 	cycles := 0
// 	maxCycles := 5

// 	for _, rune := range word {
// 		char := string(rune)
// 		position := 0
// 		lastPosition := 0
// 		positions := make([]int, 0)
// 		firstMatch := true

// 		if lastCharPosition, ok := firstCharPositions[char]; ok {
// 			position = lastCharPosition
// 			lastPosition = position
// 		} else {
// 			firstCharPositions[char] = 0
// 		}

// 		for position < maxDistance && position < maxPosition && cycles < maxCycles {
// 			// position in truncated entry
// 			position = positionDistance(char, entry[position:])

// 			// calcule real position in entry
// 			position = lastPosition + position

// 			if position >= maxDistance {
// 				// char not found, stop search
// 				break
// 			}

// 			positions = append(positions, position)

// 			// calculate next start position
// 			position = position + 1
// 			lastPosition = position

// 			if firstMatch {
// 				firstCharPositions[char] = lastPosition
// 				firstMatch = false
// 			}

// 			cycles = cycles + 1
// 		}

// 		charPositions = append(charPositions, positions)
// 	}

// 	// make a sums per index for:
// 	//    - distance of scatered characters
// 	//    - sum of all positions
// 	//    - missing character count
// 	//
// 	// characters following each other in the right order
// 	// do not increase the distance (perfect match)
// 	//
// 	// ex character positions :
// 	// [[5 15] [6 16] [7 12 20] [8 13] [9]]
// 	// ex sum by index:
// 	// [
// 	//   map[distance:5 missed:0 sum:35],
// 	//   map[distance:27 missed:0 sum:56]
// 	// ]
// 	distances := make([]map[string]int, 0)

// 	index := 0
// 	maxLength := 1

// 	foundIndexes := make([]int, 0)

// 	for index < maxLength {
// 		lastPosition := 0
// 		totalDinstance := 0
// 		totalSum := 0
// 		totalMissed := 0
// 		information := make(map[string]int)

// 		for _, positions := range charPositions {
// 			maxIndex := len(positions)

// 			if index < maxIndex {
// 				totalSum = totalSum + positions[index]

// 				if positions[index] == lastPosition+1 {
// 					lastPosition = positions[index]
// 					continue
// 				}

// 				distance := int(math.Abs(float64(positions[index] - lastPosition)))
// 				lastPosition = positions[index]

// 				if slices.Contains(foundIndexes, lastPosition) {
// 					continue
// 				}

// 				totalDinstance = totalDinstance + distance
// 			} else if index == 0 {
// 				totalMissed = totalMissed + 1
// 			}

// 			if index == 0 && maxLength < len(positions) {
// 				maxLength = len(positions)
// 			}

// 			foundIndexes = append(foundIndexes, positions...)
// 		}

// 		information["distance"] = totalDinstance
// 		information["missed"] = totalMissed
// 		information["sum"] = totalSum

// 		distances = append(distances, information)
// 		index = index + 1
// 	}

// 	if len(distances) == 0 && distances[0]["distance"] == 0 {
// 		// strings are equal
// 		// return no distance as 0
// 		return 0
// 	}

// 	// get a sum of the distances and sums
// 	// example:
// 	// [
// 	//   map[distance:5 missed:0 sum:35],
// 	//   map[distance:27 missed:0 sum:56]
// 	// ]
// 	//
// 	// sum distance 32
// 	// sum sum      91
// 	// 32 / 91 = 0.3516483516483517
// 	distance := 0
// 	sum := 0

// 	for _, information := range distances {
// 		distance = distance + information["distance"] + information["missed"]
// 		sum = sum + information["sum"]
// 	}

// 	// return distance / sum ratio
// 	return float64(distance) / float64(sum)
// }
