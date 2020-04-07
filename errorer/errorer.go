package errorer

import (
	"math/rand"
	"time"

	"github.com/mkhoi1998/mimir/consts"
)

var (
	ErrEmptyQuestion = errEmptyQuestion{}
	ErrInternal      = errInternal{}
)

type errEmptyQuestion struct{}

func (errEmptyQuestion) Error() string {
	rand.Seed(time.Now().Unix())
	return consts.Greetings[rand.Intn(len(consts.Greetings))]
}

type errInternal struct{}

func (errInternal) Error() string {
	rand.Seed(time.Now().Unix())
	return consts.ErrInternal[rand.Intn(len(consts.ErrInternal))]
}
