package iconfetcher

import (
	"strings"

	"golang.org/x/net/html"
)

func (lf *SiteIconFetcher) extractIconsFromNode(n *html.Node, base string) []iconCandidate {
	var icons []iconCandidate

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "link" {
			var rel, href, sizes, imgType string

			for _, a := range n.Attr {
				switch strings.ToLower(a.Key) {
				case "rel":
					rel = a.Val
				case "href":
					href = a.Val
				case "sizes":
					sizes = a.Val
				case "type":
					imgType = a.Val
				}
			}

			if href != "" && strings.Contains(rel, "icon") {
				iconURL := resolveURL(base, href)
				size := parseSize(sizes)
				imgType := strings.TrimPrefix(imgType, "image/")

				icons = append(icons, iconCandidate{
					url:     iconURL,
					size:    size,
					rel:     rel,
					imgType: imgType,
				})
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(n)
	return icons
}
