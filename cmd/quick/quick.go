package quick

import (
	"strings"

	"github.com/mkhoi1998/devsup/service/textrank"
	"github.com/mkhoi1998/devsup/utils"
)

func QuickChat(args []string) string {
	question := strings.Join(args, " ")
	if question == "" {
		return utils.CommentError(utils.ErrEmptyQuestion)
	}
	kw := textrank.GetKeywords(question)
	ans := strings.Join(kw, ",")
	return ans
}
