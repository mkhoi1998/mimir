package similar

import (
	"github.com/sahilm/fuzzy"
)

// GetMostSimilar return the most similar match with str from list
func GetMostSimilar(str string, list []string) string {
	matches := fuzzy.Find(str, list)
	if matches.Len() == 0 {
		return ""
	}
	return matches[0].Str
}
