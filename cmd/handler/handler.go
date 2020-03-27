package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jaytaylor/html2text"
	"golang.org/x/net/html"

	"github.com/mkhoi1998/devsup/consts"
	"github.com/mkhoi1998/devsup/service/google"
	"github.com/mkhoi1998/devsup/service/similar"
	"github.com/mkhoi1998/devsup/service/stackoverflow"
	"github.com/mkhoi1998/devsup/service/textrank"
	"github.com/mkhoi1998/devsup/service/tfidf"
	"github.com/mkhoi1998/devsup/utils"
)

// ExtractKeywords return the keywords from input question
func ExtractKeywords(args []string) []string {
	if len(args) == 0 {
		return nil
	}
	return textrank.ExtractKeywords(strings.Join(args, " "))
}

// ExtractGoogle return the content of the web page or the link gotten by Google
func ExtractGoogle(args []string) string {
	link := google.SearchGoogle(args)
	if link == "" {
		return ""
	}
	return extractContent(link)
}

// SummarizeStackWiki return the sumarized content of stack wiki
func SummarizeStackWiki(keyword string) string {
	tag, isTag := stackoverflow.CheckTagFromKeyword(keyword)
	if isTag {
		wiki, err := html2text.FromString(stackoverflow.GetWikiFromTag(tag))
		if err != nil {
			return ""
		}
		res := strings.Join(textrank.ExtractSentences(wiki, 2), "")
		return strings.ReplaceAll(res, "*", "")
	}
	return ""
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

	h := similar.GetMostSimilar(header, utils.ExtractHeaders(`(\*\*+|--+\n)?(.*)(\n(\*\*+|--+))`, wiki))
	if h == "" {
		return h
	}

	var body string
	parts := utils.ExtractBody(`((\n\n+)|(\n(\*\*+|--+)\n)|((\*\*+|--+)\n))`, wiki)
	for j := range parts {
		if parts[j] == h {
			body = parts[j+1]
			break
		}
	}

	return body
}

// SearchStackoverflow return the contents of the answer from Stackoverflow
func SearchStackoverflow(keywords []string) string {
	ans, link := stackoverflow.GetAnswerFromSearch(keywords)
	if ans == "" {
		return ans
	}
	return ParseContent(ans, link)
}

// ParseContent return the summarized content or the link if the content is too long
func ParseContent(content string, link string) string {
	isSt := link != ""
	if isSt {
		content = strings.Split(content, "<hr")[0]
	}

	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return ""
	}
	codes, err := utils.ExtractTag(doc, "code")

	// content does not have tag code
	if err != nil {
		res, err := html2text.FromString(content)
		if err != nil {
			return ""
		}

		// remove quote symbols for clearer display
		res = strings.TrimSpace(strings.ReplaceAll(res, "\n>", ""))

		// if body is too long then just return the link
		if len(strings.Split(res, "\n\n")) > 5 {
			return consts.ParseLink(link)
		}
		return res
	}
	pCodes := tfidf.GetMostImportant(codes, isSt)
	if pCodes == nil {
		return ""
	}
	return parseResponseByCode(pCodes, content, link)
}

func extractContent(link string) string {
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

		return ParseContent(ans, link)
	}

	content := google.GetContent(link)
	// do not summarize when content is too long
	if len(strings.Split(content, "\n\n")) < 17 {
		if strings.Contains(content, "<code>") {
			content := ParseContent(content, "")
			if content == "" {
				consts.ParseLink(link)
			}
			return content
		}

		return fmt.Sprintf("%v\n%v", utils.ExtractLongestBody(`(\*\*+)|(--+)`, content), consts.ParseLink(link))
	}

	return consts.ParseLink(link)
}

func parseResponseByCode(codes []string, content string, link string) string {
	isSt := link != ""

	// split parts from content
	var parts []string
	content = html.UnescapeString(content)
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

	for j := range codes {
		for i := range parts {
			if strings.Contains(parts[i], strings.Split(codes[j], "\n")[0]) {
				// ignore when code in headers
				for k := range hList {
					if strings.Contains(k, codes[j]) {
						continue
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

				r := fmt.Sprintf("\033[2;33m%v\033[0m", header)
				if index == 0 {
					r = fmt.Sprintf("%v\n\n%v", r, codes[j])
				}
				res = append(res, r)
			}
		}
	}

	if isSt {
		return fmt.Sprintf("%v\n\n%v", strings.Join(res, "\n\n"), consts.ParseLink(link))
	}

	return strings.Join(res, "\n\n")
}
