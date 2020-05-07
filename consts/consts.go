package consts

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	StackOverflowWikiBodyFilter   = "!--fGggaXKPAj"
	StackOverflowAnswerBodyFilter = "!-*jbN.9m(dML"
	StackOverflowKey              = "U4DMV*8nvpm3EOpvf69Rxw(("

	DebugHelp = "\033[1;36mJust try to explain your code to me or tell me what to do\n:h\033[0m\tI will show you how to speak to me\n\033[1;36m:q\033[0m\tGood bye.\n"
)

var (
	// Greetings store lines for greeting
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

	// Helps store lines for providing helps
	Helps = []string{
		"Hope this helps!",
	}

	// Errs store lines for response errors
	Errs = []string{
		"My mind is not responding...",
		"I cannot proceed.",
		"You may want to Google it or Stackoverflow it yourself",
	}

	// DebugGreetings store lines for greeting in debug mode
	DebugGreetings = []string{
		"Go on, I'm listening.",
		"Can you tell me what you are doing?",
	}

	// Byes store lines for farewell response
	Byes = []string{
		"I'll see you later.",
		"Happy coding :)",
	}
)

// Greet return lines from Greetings randomly
func Greet() string {
	rand.Seed(time.Now().Unix())
	return Greetings[rand.Intn(len(Greetings))]
}

// Error return lines from Errs randomly
func Error() string {
	rand.Seed(time.Now().Unix())
	return Errs[rand.Intn(len(Errs))]
}

// DebugGreet return lines from DebugGreetings randomly
func DebugGreet() string {
	rand.Seed(time.Now().Unix())
	return DebugGreetings[rand.Intn(len(DebugGreetings))]
}

// Bye return lines from Byes randomly
func Bye() string {
	rand.Seed(time.Now().Unix())
	return Byes[rand.Intn(len(Byes))]
}

// ParseLink return link with lines from Helps randomly
func ParseLink(link string) string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("%v\n%v", link, Helps[rand.Intn(len(Helps))])
}
