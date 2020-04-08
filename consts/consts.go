package consts

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	// StackOverflow stuff
	StackOverflowWikiBodyFilter   = "!--fGggaXKPAj"
	StackOverflowAnswerBodyFilter = "!-*jbN.9m(dML"
	StackOverflowKey              = "U4DMV*8nvpm3EOpvf69Rxw(("

	// Greetings store scripts for greeting
	Greetings = []string{
		"Hello world!",
		"Hello world! How can I help you?",
		"Hello world! What can I do for you today?",
		"What can I do for you today?",
		"You can ask me anything.",
		"How can I help you?",
		"How is your code?",
		"How is your code doing?",
	}

	// Helps store scripts for providing helps
	Helps = []string{
		"Hope this helps!",
	}

	// ErrInternal store scripts for response errors
	ErrInternal = []string{
		"My mind is not responding...",
		"I cannot proceed.",
	}
)

func ParseLink(link string) string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("\n\n%v\n%v", link, Helps[rand.Intn(len(Helps))])
}
