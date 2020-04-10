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

// ExtractTag return the content of the provided html tag
func ExtractTag(content, tag string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return nil, err
	}
	
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

// ExtractLongestBody return the longest body of the html
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

// ExtractHeaders return the headers of the html
func ExtractHeaders(content string) []string {
	headers := regexp.MustCompile(`(\*\*+|--+\n)?(.*)(\n(\*\*+|--+))`)
	s := headers.FindAllStringSubmatch(content, -1)
	var res []string
	for i := range s {
		res = append(res, s[i][2])
	}
	return res
}

// ExtractBody return the bodies of the html
func ExtractBody(content string) []string {
	header := regexp.MustCompile(`((\n\n+)|(\n(\*\*+|--+)\n)|((\*\*+|--+)\n))`)
	return header.Split(content, -1)
}


// RemoveAllTag remove leftover tag by html2text
func RemoveAllTag(content string) string {
	h := regexp.MustCompile(`<.*>`)
	return h.ReplaceAllString(content, "")
}

// TrimSpace remove +3 \n
func TrimSpace(content string) string {
	h := regexp.MustCompile(`(\n[\t ]*)(\n[\t ]*)(\n[\t ]*)+`)
	return h.ReplaceAllString(content, "\n\n")
}
