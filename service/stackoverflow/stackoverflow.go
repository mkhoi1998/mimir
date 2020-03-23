package stackoverflow

import (
	"github.com/mkhoi1998/Stack-on-Go/stackongo"

	"github.com/mkhoi1998/devsup/consts"
)

// CheckTagFromKeyword check if keyword exist and convert synonyms to keywords
func CheckTagFromKeyword(kw string) (string, bool) {
	session := stackongo.NewSession("stackoverflow")

	params := make(stackongo.Params)
	tag, err := session.TagInfo([]string{kw}, params)
	if err != nil {
		return kw, false
	}

	if len(tag.Items) != 0 {
		return tag.Items[0].Name, true
	}
	return kw, false
}

// GetWikiFromTag get wiki content from stackoverflow tag
func GetWikiFromTag(tag string) string {
	session := stackongo.NewSession("stackoverflow")

	params := make(stackongo.Params)
	params.Add("filter", consts.StackOverflowWikiBodyFilter)
	wiki, err := session.WikisForTags([]string{tag}, params)
	if err != nil {
		return ""
	}

	if len(wiki.Items) != 0 {
		return wiki.Items[0].Body
	}
	return ""
}

// GetAnswerFromSearch from question with > 99 votes and accepted or highest-vote answer
func GetAnswerFromSearch(query []string) string {
	session := stackongo.NewSession("stackoverflow")

	params := make(stackongo.Params)
	params.Sort("relevance")
	items, err := session.AdvancedSearch(query, params)
	if err != nil {
		return ""
	}

	if len(items.Items) != 0 && items.Items[0].Score > 99 {
		if items.Items[0].Accepted_answer_id != 0 {
			params = make(stackongo.Params)
			params.Add("filter", consts.StackOverflowAnswerBodyFilter)
			ans, err := session.GetAnswers([]int{items.Items[0].Accepted_answer_id}, params)
			if err != nil {
				return ""
			}

			if len(ans.Items) != 0 {
				return ans.Items[0].Body
			}
		}

		return GetAnswerFromQuestionID(items.Items[0].Question_id)
	}
	return ""
}

// GetAnswerFromQuestionID return most voted answer from question id
func GetAnswerFromQuestionID(id int) string {
	session := stackongo.NewSession("stackoverflow")

	params := make(stackongo.Params)
	params.Add("filter", consts.StackOverflowAnswerBodyFilter)
	params.Sort("votes")
	ans, err := session.AnswersForQuestions([]int{id}, params)
	if err != nil {
		return ""
	}

	if len(ans.Items) != 0 {
		return ans.Items[0].Body
	}
	return ""
}
