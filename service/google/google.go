package google

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/jaytaylor/html2text"
	"github.com/mkhoi1998/devsup/utils"
)

// SearchGoogle get the first link from google search using queries
func SearchGoogle(query []string) string {
	u := fmt.Sprintf("https://www.google.com/search?q=%v", url.QueryEscape(strings.Join(query, " ")))
	r, err := http.Get(u)
	if err != nil {
		return ""
	}

	t, err := html2text.FromReader(r.Body)
	if err != nil {
		return ""
	}

	rawText := strings.ReplaceAll(utils.ExtractLongestBody(`\*\*+`, t), "~", ".")

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
		return ""
	}

	return strings.TrimPrefix(res, "q=")
}

// GetContent from the link and parse it to text
func GetContent(link string) string {
	r, err := http.Get(link)
	if err != nil {
		return ""
	}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	b := string(bodyBytes)
	if strings.Contains(b, "<code>") {
		return b
	}
	t, err := html2text.FromString(b)
	if err != nil {
		return ""
	}

	return t
}
