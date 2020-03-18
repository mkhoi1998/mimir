package utils

import (
	"errors"
)

var (
	ErrEmptyQuestion = errors.New("empty_question")
)

func CommentError(err error) string {
	switch err {
	case ErrEmptyQuestion:
		return "You can ask me anything"
	}
	return ""
}
