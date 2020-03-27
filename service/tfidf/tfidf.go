package tfidf

import (
	"sort"
	"strings"

	"github.com/wilcosheh/tfidf"
)

// GetMostImportant returns the paragraphs in order of importances desc
func GetMostImportant(text []string, isSt bool) []string {
	type weightData struct {
		Data   string
		Weight float64
	}
	f := tfidf.New()
	var temp []string
	for i := range text {
		temp = append(temp, text[i])
	}
	for i := range temp {
		temp[i] = strings.ReplaceAll(temp[i], "[", " ")
		temp[i] = strings.ReplaceAll(temp[i], "]", " ")
		temp[i] = strings.ReplaceAll(temp[i], "(", " ")
		temp[i] = strings.ReplaceAll(temp[i], ")", " ")
		temp[i] = strings.ReplaceAll(temp[i], ";", " ")
		temp[i] = strings.ReplaceAll(temp[i], ".", " ")
		temp[i] = strings.ToLower(temp[i])
	}
	f.AddDocs(temp...)
	var wD []weightData
	for i := range temp {
		if len(strings.Fields(text[i])) < 2 {
			continue
		}
		w := f.Cal(temp[i])
		var weight float64
		for j := range w {
			weight += w[j]
		}
		if isSt && len(wD) != 0 {
			if weight > 1 {
				continue
			}
		}
		weight /= float64(len(w))
		wD = append(wD, weightData{Weight: weight, Data: text[i]})
	}
	sort.Slice(wD[:], func(i, j int) bool {
		return wD[i].Weight < wD[j].Weight
	})
	if len(wD) != 0 && len(wD) != 1 {
		if wD[0].Weight == wD[len(wD)-1].Weight {
			return nil
		}
	}

	var res []string
	for i := range wD {
		res = append(res, wD[i].Data)
	}
	if len(res) > 3 {
		res = res[:3]
	}
	return res
}
