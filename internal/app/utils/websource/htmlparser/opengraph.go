package htmlparser

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

func ExtractOpenGraphImageURLs(body []byte) ([]string, error) {
	var imageURLs []string

	if len(body) == 0 {
		return nil, nil
	}

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var head *html.Node
	for node := doc.FirstChild; node != nil; node = node.NextSibling {
		if node.Type != html.ElementNode || node.Data != "html" {
			continue
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode && child.Data == "head" {
				head = child
				break
			}
		}

		break
	}

	if head == nil {
		return nil, nil
	}

	for node := head.FirstChild; node != nil; node = node.NextSibling {
		if node.Type != html.ElementNode || node.Data != "meta" {
			continue
		}

		var property, content string

		for _, attribute := range node.Attr {
			switch attribute.Key {
			case "property":
				property = attribute.Val
			case "content":
				content = attribute.Val
			}
		}

		if strings.EqualFold(property, "og:image") && content != "" {
			imageURLs = append(imageURLs, content)
		}
	}

	return imageURLs, nil
}
