package logofetcher

import (
	"strings"

	"golang.org/x/net/html"
)

func (lf *SiteLogoFetcher) extractIcons(links [][]html.Attribute, base string) []iconCandidate {
	var icons []iconCandidate

	for _, link := range links {
		var rel, href, sizes, imgType string
		for _, a := range link {
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

	return icons
}
