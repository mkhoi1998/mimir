package errorer

import (
	"math/rand"
	"time"

	"github.com/mkhoi1998/devsup/consts"
)

var (
	ErrEmptyQuestion = errEmptyQuestion{}
)

type errEmptyQuestion struct{}

func (errEmptyQuestion) Error() string {
	rand.Seed(time.Now().Unix())
	return consts.Greetings[rand.Intn(len(consts.Greetings))]
}
