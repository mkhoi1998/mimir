package textrank

import (
	textrank "github.com/DavidBelicza/TextRank"
)

func GetKeywords(text string) []string {
	tr := textrank.NewTextRank()
	rule := textrank.NewDefaultRule()
	language := textrank.NewDefaultLanguage()
	algorithmDef := textrank.NewDefaultAlgorithm()

	tr.Populate(text, language, rule)
	tr.Ranking(algorithmDef)

	w := textrank.FindSingleWords(tr)
	var res []string
	for i := range w {
		res = append(res, w[i].Word)
	}
	return res
}
