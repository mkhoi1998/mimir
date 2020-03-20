package cmd

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/jaytaylor/html2text"

	"github.com/mkhoi1998/devsup/consts"
	"github.com/mkhoi1998/devsup/errorer"
	"github.com/mkhoi1998/devsup/service/google"
	"github.com/mkhoi1998/devsup/service/stackoverflow"
	"github.com/mkhoi1998/devsup/service/textrank"
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
			res = parseContent(link)
		}
		return res

	default:
		res := searchStackOverflow(keywords)
		if res == "" {
			link := google.SearchGoogle(args)
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
	ans, err := html2text.FromString(ans)
	if err != nil {
		return errorer.ErrInternal.Error()
	}
	if len(strings.Split(ans, "\n\n")) < 20 {
		return ans
	}
	res := strings.Join(textrank.ExtractSentences(ans, 1), "")
	return strings.ReplaceAll(res, "*", "")
}

func parseContent(link string) string {
	content := google.GetContent(link)
	if len(strings.Split(content, "\n\n")) < 17 && !strings.Contains(strings.ToLower(content), "captcha") {
		header := regexp.MustCompile(`(\*\*+)|(--+)`)
		ts := header.Split(content, -1)

		var index int
		var last int
		for i := range ts {
			if len(ts[i]) > last {
				index = i
				last = len(ts[i])
			}
		}
		return fmt.Sprintf("%v\n%v", ts[index], link)
	}

	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("%v\n%v", link, consts.Helps[rand.Intn(len(consts.Helps))])
}
