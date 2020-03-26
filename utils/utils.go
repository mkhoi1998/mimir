package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func ExtractTag(doc *html.Node, tag string) ([]string, error) {
	var code []string
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == tag {
			for parent := node; parent != nil; parent = parent.Parent {
				if parent.Data == "blockquote" {
					return
				}
			}
			code = append(code, html.UnescapeString(strings.TrimSpace(renderNode(node, tag))))
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if code != nil {
		return code, nil
	}
	return nil, errors.New("Missing tag in the node tree")
}

func renderNode(n *html.Node, tag string) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return strings.ReplaceAll(strings.ReplaceAll(buf.String(), fmt.Sprintf("<%v>", tag), ""), fmt.Sprintf("</%v>", tag), "")
}

func ExtractLongestBody(regex, content string) string {
	header := regexp.MustCompile(regex)
	ts := header.Split(content, -1)

	var index int
	var last int
	for i := range ts {
		if len(ts[i]) > last {
			index = i
			last = len(ts[i])
		}
	}
	return ts[index]
}

func ExtractHeaders(regex, content string) []string {
	headers := regexp.MustCompile(regex)
	s := headers.FindAllStringSubmatch(content, -1)
	var res []string
	for i := range s {
		res = append(res, s[i][2])
	}
	return res
}

func ExtractBody(regex, content string) []string {
	header := regexp.MustCompile(regex)
	return header.Split(content, -1)
}
