package stackoverflow

import (
	"github.com/mkhoi1998/Stack-on-Go/stackongo"
)

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

func GetWikiFromTag(tag string) string {
	session := stackongo.NewSession("stackoverflow")
	params := make(stackongo.Params)
	params.Add("filter", "!--fGggaXKPAj")

	wiki, err := session.WikisForTags([]string{tag}, params)
	if err != nil {
		return ""
	}
	if len(wiki.Items) != 0 {
		return wiki.Items[0].Body
	}
	return ""
}

func GetAnswerFromSearch(query []string) string {
	session := stackongo.NewSession("stackoverflow")
	params := make(stackongo.Params)
	params.Sort("relevance")

	items, err := session.AdvancedSearch(query, params)
	if err != nil {
		return ""
	}
	if len(items.Items) != 0 && items.Items[0].Score > 10 && items.Items[0].Accepted_answer_id != 0 {
		params = make(stackongo.Params)
		params.Add("filter", "!9Z(-wzu0T")
		ans, err := session.GetAnswers([]int{items.Items[0].Accepted_answer_id}, params)
		if err != nil {
			return ""
		}
		if len(ans.Items) != 0 {
			return ans.Items[0].Body
		}
	}
	return ""
}
