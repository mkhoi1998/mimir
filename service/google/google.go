package google

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/mkhoi1998/mimir/utils"
)

// SearchGoogle return the first link from google search using queries
func SearchGoogle(query []string) string {
	r, err := http.Get(fmt.Sprintf("https://www.google.com/search?q=%v", url.QueryEscape(strings.Join(query, " "))))
	if err != nil {
		return ""
	}

	var res string

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	urlReg := regexp.MustCompile(`q=https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
	matches := urlReg.FindAllString(strings.ReplaceAll(utils.ExtractLongestBody(`\*\*+`, string(bodyBytes)), "~", "."), -1)
	for i := range matches {
		if strings.Contains(matches[i], "https://www.google.com") {
			continue
		}
		res = strings.TrimSuffix(strings.TrimPrefix(matches[i], "q="), "&amp")
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
	if !strings.Contains(r.Header.Get("content-type"), "html") {
		return ""
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	b := string(bodyBytes)
	b = html.UnescapeString(b)
	option := regexp.MustCompile(`<option.*>.*<\/option>`)
	b = option.ReplaceAllString(b, "")
	return b
}
