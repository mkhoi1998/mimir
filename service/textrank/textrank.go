package textrank

import (
	"regexp"
	"strings"

	textrank "github.com/DavidBelicza/TextRank"
)

func ExtractKeywords(text string) []string {
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

func ExtractSentences(text string, count int) []string {
	r := regexp.MustCompile(`([a-z0-9])(\.[a-z0-9])+`)
	rep := r.FindAllString(text, -1)
	for i := range rep {
		temp := strings.ReplaceAll(rep[i], ".", "~")
		text = strings.ReplaceAll(text, rep[i], temp)
	}

	tr := textrank.NewTextRank()
	rule := textrank.NewDefaultRule()
	language := textrank.NewDefaultLanguage()
	algorithmDef := textrank.NewDefaultAlgorithm()

	tr.Populate(text, language, rule)
	tr.Ranking(algorithmDef)

	s := textrank.FindSentencesByWordQtyWeight(tr, count)
	var res []string
	for i := range s {
		res = append(res, strings.ReplaceAll(s[i].Value, "~", "."))
	}
	return res
}
