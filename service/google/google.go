package google

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/jaytaylor/html2text"

	"github.com/mkhoi1998/mimir/utils"
)

// SearchGoogle return the first link from google search using queries
func SearchGoogle(query []string) string {
	r, err := http.Get(fmt.Sprintf("https://www.google.com/search?q=%v", url.QueryEscape(strings.Join(query, " "))))
	if err != nil {
		return ""
	}

	t, err := html2text.FromReader(r.Body)
	if err != nil {
		return ""
	}

	var res string

	urlReg := regexp.MustCompile(`q=https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
	matches := urlReg.FindAllString(strings.ReplaceAll(utils.ExtractLongestBody(`\*\*+`, t), "~", "."), -1)
	for i := range matches {
		if strings.Contains(matches[i], "https://www.google.com") {
			continue
		}
		res = strings.Split(strings.TrimPrefix(matches[i], "q="), "&sa=")[0]
		break
	}
	res, err = url.QueryUnescape(res)
	if err != nil {
		return ""
	}

	return res
}

// GetContent return html content from link and parse to text
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

	t, err := html2text.FromString(b, html2text.Options{OmitLinks: true})
	if err != nil {
		return ""
	}
	return t
}
