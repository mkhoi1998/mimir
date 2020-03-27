package answer

import (
	"fmt"

	"github.com/mkhoi1998/devsup/cmd/handler"
	"github.com/mkhoi1998/devsup/errorer"
)

// Handler response input question based on keywords and online resources (Stackoverflow and Google)
func Handler(args []string) {
	fmt.Println(responseAnswer(args))
}

func responseAnswer(args []string) string {
	keywords := handler.ExtractKeywords(args)

	switch len(keywords) {
	case 0:
		return errorer.ErrEmptyQuestion.Error()

	case 1:
		res := handler.SummarizeStackWiki(keywords[0])
		if res == "" {
			res = handler.ExtractGoogle(args)
		}
		if res == "" {
			res = errorer.ErrInternal.Error()
		}
		return res

	case 2:
		res := handler.ExtractStackWiki(keywords)
		if res == "" {
			res = handler.ExtractGoogle(args)
		}
		if res == "" {
			res = errorer.ErrInternal.Error()
		}
		return res

	default:
		res := handler.SearchStackoverflow(keywords)
		if res == "" {
			res = handler.ExtractGoogle(args)
		}
		if res == "" {
			res = errorer.ErrInternal.Error()
		}
		return res
	}
}
