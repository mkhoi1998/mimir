package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jaytaylor/html2text"

	"github.com/mkhoi1998/mimir/consts"
	"github.com/mkhoi1998/mimir/service/google"
	"github.com/mkhoi1998/mimir/service/similar"
	"github.com/mkhoi1998/mimir/service/stackoverflow"
	"github.com/mkhoi1998/mimir/service/textrank"
	"github.com/mkhoi1998/mimir/service/tfidf"
	"github.com/mkhoi1998/mimir/utils"
)

// ExtractKeywords return the keywords from input question
func ExtractKeywords(args []string) []string {
	if len(args) == 0 {
		return nil
	}
	return textrank.ExtractKeywords(strings.Join(args, " "))
}

// SummarizeStackWiki return the sumarized content of stack wiki
func SummarizeStackWiki(keyword string) string {
	tag, isTag := stackoverflow.CheckTagFromKeyword(keyword)
	if !isTag {
		return ""
	}

	wiki, err := html2text.FromString(stackoverflow.GetWikiFromTag(tag))
	if err != nil {
		return ""
	}
	res := strings.Join(textrank.ExtractSentences(wiki, 2), "")
	return strings.ReplaceAll(res, "*", "")
}

// ExtractStackWiki return body of the specified header of the stack wiki
func ExtractStackWiki(keywords []string) string {
	tag, isTag := stackoverflow.CheckTagFromKeyword(keywords[0])
	header := keywords[1]
	if !isTag {
		tag, isTag = stackoverflow.CheckTagFromKeyword(keywords[1])
		if !isTag {
			return ""
		}
		header = keywords[0]
	}

	wiki, err := html2text.FromString(stackoverflow.GetWikiFromTag(tag))
	if err != nil {
		return ""
	}

	h := similar.GetMostSimilar(header, utils.ExtractHeaders(wiki))
	if h == "" {
		return h
	}

	var body string
	parts := utils.ExtractBody(wiki)
	for j := range parts {
		if parts[j] == h {
			body = parts[j+1]
			break
		}
	}

	return body
}

// SummarizeStackoverflow return the summarized contents of the answer from Stackoverflow
func SummarizeStackoverflow(keywords []string) string {
	ans, link := stackoverflow.GetAnswerFromSearch(keywords)
	if ans == "" {
		return ans
	}
	return sumarizeSmallContent(strings.Split(ans, "<hr")[0], link)
}

// SummarizeGoogle return the summarized content of the web page or the link gotten by Google
func SummarizeGoogle(args []string) string {
	link := google.SearchGoogle(args)
	if link == "" {
		return ""
	}

	if strings.Contains(link, "stackoverflow.com") {
		s := strings.Split(link, "/")
		id, err := strconv.Atoi(s[4])
		if err != nil {
			return ""
		}

		ans, _ := stackoverflow.GetAnswerFromQuestionID(id)
		if ans == "" {
			return consts.ParseLink(link)
		}

		return sumarizeSmallContent(strings.Split(ans, "<hr")[0], link)
	}

	content := google.GetContent(link)
	if len(strings.Split(content, "\n\n")) < 17 {
		if strings.Contains(content, "<code>") {
			content := sumarizeSmallContent(content, "")
			if content == "" {
				consts.ParseLink(link)
			}
			return fmt.Sprintf("%v%v", content, consts.ParseLink(link))
		}

		return fmt.Sprintf("%v\n%v", utils.ExtractLongestBody(`(\*\*+)|(--+)`, content), consts.ParseLink(link))
	}

	ct := utils.ExtractBody(content)
	var cts []string
	dup := map[string]bool{}
	for i := range ct {
		if ct[i] == "" {
			continue
		}
		temp := strings.Split(ct[i], "\n")
		isContent := true
		for j := range temp {
			if temp[j] == "" {
				continue
			}
			if len(strings.Split(temp[j], " ")) < 10 || len(temp[j]) < 100 {
				if len(ct[i]) > 400 && !strings.Contains(ct[i], "*") &&
					!strings.Contains(ct[i], ".jpg") && !strings.Contains(ct[i], ".jpeg") &&
					!strings.Contains(ct[i], ".png") && !strings.Contains(ct[i], ".gif") {
					continue
				}
				isContent = false
				break
			}
		}
		if isContent {
			if !dup[ct[i]] {
				cts = append(cts, utils.RemoveAllTag(ct[i]))
				dup[ct[i]] = true
			}
			if i < len(ct)-2 {
				if !dup[ct[i+1]] {
					cts = append(cts, utils.RemoveAllTag(ct[i+1]))
					dup[ct[i+1]] = true
				}
			}

		}
	}
	for i := range cts {
		cts[i] = strings.ReplaceAll(cts[i], "*", "")
	}
	var res string
	if len(cts) > 20 {
		res = fmt.Sprintf("%v%v", strings.Join(textrank.ExtractSentences(strings.Join(cts, "\n"), 1), "\n\n"), consts.ParseLink(link))
	} else {
		res = fmt.Sprintf("%v%v", strings.Join(cts, "\n\n"), consts.ParseLink(link))
	}
	if strings.Count(res, "\n") >= 50 {
		return consts.ParseLink(link)
	}

	return utils.TrimSpace(res)
}

// GetStackWiki return the body of the stack wiki
func GetStackWiki(keyword string) []string {
	tag, isTag := stackoverflow.CheckTagFromKeyword(keyword)
	if !isTag {
		return nil
	}

	wiki, err := html2text.FromString(stackoverflow.GetWikiFromTag(tag))
	if err != nil {
		return nil
	}

	return utils.ExtractBody(wiki)
}

// GetStackoverflow return the full contents of the answer from Stackoverflow
func GetStackoverflow(keywords []string) []string {
	ans, link := stackoverflow.GetAnswerFromSearch(keywords)
	if ans == "" {
		return nil
	}

	res, err := html2text.FromString(ans)
	if err != nil {
		return nil
	}

	// remove quote symbols for clearer display
	res = strings.TrimSpace(strings.ReplaceAll(res, "\n>", ""))

	return append(strings.Split(res, "\n\n"), consts.ParseLink(link))
}

func sumarizeSmallContent(content, link string) string {
	codes, err := utils.ExtractTag(content, "code")
	if err == nil {
		return summarizeCodeContent(codes, content, link)
	}
	// content does not have tag code
	res, err := html2text.FromString(content)
	if err != nil {
		return ""
	}

	// remove quote symbols for clearer display
	res = strings.TrimSpace(strings.ReplaceAll(res, "\n>", ""))

	// if body is too long then just return the link
	if len(strings.Split(res, "\n\n")) > 5 {
		if link != "" {
			return consts.ParseLink(link)
		}
		return link
	}
	return res

}

func summarizeCodeContent(codes []string, content string, link string) string {
	isSt := link != ""

	// split parts from content
	var parts []string
	if isSt {
		parts = strings.Split(content, "\n\n")
	} else {
		content = strings.ReplaceAll(content, "\n\n", "\n")
		parts = strings.Split(content, "\n")
	}
	for i := range parts {
		parts[i] = strings.ReplaceAll(parts[i], "\n", " ")
	}

	// pList to check duplicate headers for code
	hList := map[string]bool{}
	var res []string

	codes = tfidf.GetMostImportant(codes, isSt)
	if codes == nil {
		return ""
	}

	for j := range codes {
		for i := range parts {
			if strings.Contains(parts[i], strings.Split(codes[j], "\n")[0]) {
				// ignore when code in headers
				isDuplicate := false
				for k := range hList {
					if strings.Contains(k, codes[j]) {
						isDuplicate = true
					}
				}

				index := i - 1
				if i == 0 {
					index = 0
				}

				header, err := html2text.FromString(parts[index])
				if err != nil {
					return ""
				}

				// ignore duplicate headers
				if hList[header] == true {
					continue
				}
				hList[header] = true

				r := fmt.Sprintf("\033[0;36m%v\033[0m", header)
				if !isDuplicate {
					r = fmt.Sprintf("%v\n\n%v", r, codes[j])
				}
				res = append(res, r)
			}
		}
	}

	if isSt {
		return fmt.Sprintf("%v%v", strings.Join(res, "\n\n"), consts.ParseLink(link))
	}

	return strings.Join(res, "\n\n")
}
