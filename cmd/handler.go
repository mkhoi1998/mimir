package cmd

import (
	"strings"

	"github.com/jaytaylor/html2text"
	"github.com/mkhoi1998/devsup/errorer"
	"github.com/mkhoi1998/devsup/service/stackoverflow"
	"github.com/mkhoi1998/devsup/service/textrank"
)

func ResponseHandler(args []string) string {
	question := strings.Join(args, " ")
	if question == "" {
		return errorer.ErrEmptyQuestion.Error()
	}
	keywords := textrank.ExtractKeywords(question)

	switch len(keywords) {
	case 1:
		res := tagStackOverflow(keywords[0])
		return res

	default:
		return searchStackOverflow(keywords)
	}
}

func tagStackOverflow(keyword string) string {
	tag, isTag := stackoverflow.CheckTagFromKeyword(keyword)
	if isTag {
		wiki := stackoverflow.GetWikiFromTag(tag)
		wiki, err := html2text.FromString(wiki)
		if err != nil {
			return errorer.ErrEmptyQuestion.Error()
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
		return errorer.ErrEmptyQuestion.Error()
	}
	if len(strings.Split(ans, "\n\n")) < 20 {
		return ans
	}
	res := strings.Join(textrank.ExtractSentences(ans, 1), "")
	return strings.ReplaceAll(res, "*", "")
}
