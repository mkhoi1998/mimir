package cmd

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/jaytaylor/html2text"
	"golang.org/x/net/html"

	"github.com/mkhoi1998/devsup/consts"
	"github.com/mkhoi1998/devsup/errorer"
	"github.com/mkhoi1998/devsup/service/google"
	"github.com/mkhoi1998/devsup/service/similar"
	"github.com/mkhoi1998/devsup/service/stackoverflow"
	"github.com/mkhoi1998/devsup/service/textrank"
	"github.com/mkhoi1998/devsup/service/tfidf"
	"github.com/mkhoi1998/devsup/utils"
)

// ResponseHandler return answer based on user question
func ResponseHandler(args []string) string {
	question := strings.Join(args, " ")
	if question == "" {
		return errorer.ErrEmptyQuestion.Error()
	}
	keywords := textrank.ExtractKeywords(question)
	switch len(keywords) {
	case 1:
		res := tagStackOverflow(keywords[0])
		if res == "" {
			link := google.SearchGoogle(args)
			if link == "" {
				return errorer.ErrInternal.Error()
			}
			res = parseContent(link)
		}
		return res

	case 2:
		tag, isTag := stackoverflow.CheckTagFromKeyword(keywords[0])
		other := keywords[1]
		if !isTag {
			tag, isTag = stackoverflow.CheckTagFromKeyword(keywords[1])
			if !isTag {
				res := searchStackOverflow(keywords)
				if res == "" {
					link := google.SearchGoogle(args)
					if link == "" {
						return errorer.ErrInternal.Error()
					}
					res = parseContent(link)
				}
				return res
			}
			other = keywords[0]
		}
		wiki := stackoverflow.GetWikiFromTag(tag)
		wiki, err := html2text.FromString(wiki)
		if err != nil {
			return errorer.ErrInternal.Error()
		}
		h := similar.GetMostSimilar(other, utils.ExtractHeaders(`(\*\*+|--+\n)?(.*)(\n(\*\*+|--+))`, wiki))
		if h == "" {
			link := google.SearchGoogle(args)
			if link == "" {
				return errorer.ErrInternal.Error()
			}
			return parseContent(link)
		}
		parts := utils.ExtractBody(`((\n\n+)|(\n(\*\*+|--+)\n)|((\*\*+|--+)\n))`, wiki)
		var body string
		for j := range parts {
			if parts[j] == h {
				body = parts[j+1]
				break
			}
		}

		return body

	default:
		res := searchStackOverflow(keywords)
		if res == "" {
			link := google.SearchGoogle(args)
			if link == "" {
				return errorer.ErrInternal.Error()
			}
			res = parseContent(link)
		}
		return res
	}
}

func tagStackOverflow(keyword string) string {
	tag, isTag := stackoverflow.CheckTagFromKeyword(keyword)
	if isTag {
		wiki := stackoverflow.GetWikiFromTag(tag)
		wiki, err := html2text.FromString(wiki)
		if err != nil {
			return errorer.ErrInternal.Error()
		}
		res := strings.Join(textrank.ExtractSentences(wiki, 2), "")
		return strings.ReplaceAll(res, "*", "")
	}
	return ""
}

func searchStackOverflow(keywords []string) string {
	ans := stackoverflow.GetAnswerFromSearch(keywords)
	if ans == "" {
		return ans
	}
	return extractContent(ans, true)
}

func parseContent(link string) string {
	if strings.Contains(link, "stackoverflow.com") {
		s := strings.Split(link, "/")
		id, err := strconv.Atoi(s[4])
		if err != nil {
			return errorer.ErrInternal.Error()
		}

		ans := stackoverflow.GetAnswerFromQuestionID(id)
		if ans == "" {
			return fmt.Sprintf("%v\n%v", link, consts.Helps[rand.Intn(len(consts.Helps))])
		}
		return extractContent(ans, true)
	}

	content := google.GetContent(link)
	if len(strings.Split(content, "\n\n")) < 17 {
		if strings.Contains(content, "<code>") {
			c := extractContent(content, false)
			if c == "" {
				rand.Seed(time.Now().Unix())
				return fmt.Sprintf("%v\n%v", link, consts.Helps[rand.Intn(len(consts.Helps))])
			}
			return c
		}

		return fmt.Sprintf("%v\n%v", utils.ExtractLongestBody(`(\*\*+)|(--+)`, content), link)
	}

	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("%v\n%v", link, consts.Helps[rand.Intn(len(consts.Helps))])
}

func extractContent(content string, isSt bool) string {
	if isSt {
		content = strings.Split(content, "<hr")[0]
	}
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return errorer.ErrInternal.Error()
	}
	codes, err := utils.ExtractTag(doc, "code")
	if err != nil {
		res, err := html2text.FromString(content)
		if err != nil {
			return errorer.ErrInternal.Error()
		}
		return res
	}
	extracted := tfidf.GetMostImportant(codes, isSt)
	if extracted == nil {
		return ""
	}
	return parseResponseByCode(extracted, content, isSt)
}

func parseResponseByCode(codes []string, ans string, isSt bool) string {
	ans = html.UnescapeString(ans)
	var parts []string
	if isSt {
		parts = strings.Split(ans, "\n\n")
	} else {
		ans = strings.ReplaceAll(ans, "\n\n", "\n")
		parts = strings.Split(ans, "\n")
	}

	pList := map[string]bool{}
	for j := range codes {
		for i := range parts {
			parts[i] = strings.ReplaceAll(parts[i], "\n", " ")
			if strings.Contains(parts[i], strings.Split(codes[j], "\n")[0]) {
				for k := range pList {
					if strings.Contains(k, codes[j]) {
						codes[j] = ""
						continue
					}
				}
				index := i - 1
				if i == 0 {
					index = 0
					codes[j] = ""
				}
				prefix, err := html2text.FromString(parts[index])
				if err != nil {
					return errorer.ErrInternal.Error()
				}

				if pList[prefix] == true {
					continue
				}
				pList[prefix] = true
				if codes[j] == "" {
					codes[j] = fmt.Sprintf("\033[2;33m%s\033[0m", prefix)
				} else {
					codes[j] = fmt.Sprintf("\033[2;33m%s\033[0m\n\n%v", prefix, codes[j])
				}
			}
		}
	}
	return strings.Join(codes, "\n\n")
}
