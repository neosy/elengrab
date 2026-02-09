package logofetcher

import (
	"context"
	"errors"
	"net/url"
	"sort"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	"github.com/neosy/elengrab/pkg/httpx"
)

// FetchLogo fetches the logo from the given site and returns ImageData.
func (lf *SiteLogoFetcher) FetchLogo(ctx context.Context, siteURL string) (*dmedia.ImageData, error) {
	links, err := httpx.GetLinksInHead(
		ctx,
		siteURL,
		httpx.GetOptions{
			Limit:            limitHTML,
			IgnoreStatusCode: true,
		},
		httpx.ClientOptionWithTimeout(getHTMLTimeout),
	)
	if err != nil {
		return nil, err
	}

	candidates := lf.extractIcons(links, siteURL)

	// fallback
	if len(candidates) == 0 {
		u, _ := url.Parse(siteURL)
		u.Path = "/favicon.ico"
		candidates = append(candidates, iconCandidate{url: u.String()})
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		iIsSVG := candidates[i].imgType == "svg+xml"
		jIsSVG := candidates[j].imgType == "svg+xml"

		// SVG has higher priority
		if iIsSVG && !jIsSVG {
			return true
		}

		// if i is not SVG and j is SVG → i < j → false
		if iIsSVG || jIsSVG {
			return false
		}

		// Both not SVG → prefer larger size
		return candidates[i].size > candidates[j].size
	})

	for _, c := range candidates {
		logo, err := lf.downloadImage(ctx, c.url)
		if err == nil {
			return logo, nil
		}
	}

	return nil, errors.New("logo not found")
}
