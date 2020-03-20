package google

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/jaytaylor/html2text"

	"github.com/mkhoi1998/devsup/errorer"
)

// SearchGoogle get the first link from google search using queries
func SearchGoogle(query []string) string {
	u := fmt.Sprintf("https://www.google.com/search?q=%v", url.QueryEscape(strings.Join(query, " ")))
	r, err := http.Get(u)
	if err != nil {
		return errorer.ErrInternal.Error()
	}

	t, err := html2text.FromReader(r.Body)
	if err != nil {
		return errorer.ErrInternal.Error()
	}

	header := regexp.MustCompile(`\*\*+`)
	ts := header.Split(t, -1)

	var index int
	var last int
	for i := range ts {
		if len(ts[i]) > last {
			index = i
			last = len(ts[i])
		}
	}
	rawText := strings.ReplaceAll(ts[index], "~", ".")

	urlReg := regexp.MustCompile(`q=https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
	matches := urlReg.FindAllString(rawText, -1)
	var res string
	for i := range matches {
		if strings.Contains(matches[i], "https://www.google.com") {
			continue
		}
		res = matches[i]
		break
	}
	res, err = url.QueryUnescape(strings.Split(res, "&sa=")[0])
	if err != nil {
		return errorer.ErrInternal.Error()
	}

	return strings.TrimPrefix(res, "q=")
}

// GetContent from the link and parse it to text
func GetContent(link string) string {
	r, err := http.Get(link)
	if err != nil {
		return errorer.ErrInternal.Error()
	}
	t, err := html2text.FromReader(r.Body)
	if err != nil {
		return errorer.ErrInternal.Error()
	}

	return t
}
