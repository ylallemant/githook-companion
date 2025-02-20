package commit

import (
	edlib "github.com/hbollon/go-edlib"
	"github.com/ylallemant/githook-companion/pkg/api"
)

func fuzzyDictionaryMatch(token string, config *api.Config) *api.CommitTypeDictionary {
	var match *api.CommitTypeDictionary
	minDistance := 100000

	for _, dictionary := range config.Commit.Dictionaries {
		for _, synonym := range dictionary.Synonyms {
			distance := edlib.DamerauLevenshteinDistance(token, synonym)

			//fmt.Println("   -", dictionary.Name, token, synonym, distance)
			if distance < minDistance {
				minDistance = distance
				match = dictionary
			}
		}
	}

	if minDistance > 2 {
		return nil
	}

	return match
}
