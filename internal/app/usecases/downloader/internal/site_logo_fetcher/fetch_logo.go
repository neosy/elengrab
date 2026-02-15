package logofetcher

import (
	"context"
	"errors"
	"net/url"
	"sort"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	"github.com/neosy/elengrab/pkg/httpx"
)

var ErrLogoNotFound = errors.New("logo not found")

// FetchBestLogo fetches the best available logo for the given site URL.
// It retrieves icon candidates from the site, downloads them in priority order,
// and returns the first successfully downloaded image.
// Returns ErrLogoNotFound if no valid logo could be retrieved.
func (lf *SiteLogoFetcher) FetchBestLogo(ctx context.Context, siteURL string) (*dmedia.ImageData, error) {
	candidates, err := lf.fetchCandidatesFromURL(ctx, siteURL)
	if err != nil {
		return nil, err
	}

	for _, c := range candidates {
		image, err := lf.downloadImage(ctx, c.url)
		if err == nil {
			return image, nil
		}
	}

	return nil, ErrLogoNotFound
}

// FetchLogos fetches all available logos for the given site URL.
// It retrieves icon candidates from the site, downloads each candidate,
// and returns all successfully downloaded images.
func (lf *SiteLogoFetcher) FetchLogos(ctx context.Context, siteURL string) ([]*dmedia.ImageData, error) {
	candidates, err := lf.fetchCandidatesFromURL(ctx, siteURL)
	if err != nil {
		return nil, err
	}

	capHint := min(len(candidates), 4)
	logos := make([]*dmedia.ImageData, 0, capHint)
	for _, c := range candidates {
		logo, err := lf.downloadImage(ctx, c.url)
		if err == nil {
			logos = append(logos, logo)
		}
	}

	return logos, nil
}

// fetchCandidatesFromURL retrieves icon candidates from the site's HTML head,
// applies a fallback to /favicon.ico, and sorts candidates by priority.
func (lf *SiteLogoFetcher) fetchCandidatesFromURL(ctx context.Context, siteURL string) ([]iconCandidate, error) {
	links, err := httpx.GetLinksInHead(
		ctx,
		siteURL,
		httpx.GetOptions{
			Limit:            limitHTML,
			IgnoreStatusCode: true,
		},
		httpx.ClientOptionWithTimeout(fetchLogoTimeout),
		httpx.ClientOptionWithDefaultCookieJar(),
	)
	if err != nil {
		return nil, err
	}

	var candidates []iconCandidate
	if len(links) > 0 {
		candidates = lf.extractIcons(links, siteURL)
	} else {
		// fallback to default favicon.ico if no icons were found
		u, err := url.Parse(siteURL)
		if err == nil {
			u.Path = "/favicon.ico"
			candidates = append(candidates, iconCandidate{url: u.String()})
		}
	}

	if len(candidates) == 0 {
		return candidates, nil
	}

	// sort candidates by preference:
	// 1. SVG icons have the highest priority
	// 2. Non-SVG icons are sorted by descending size
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

	return candidates, nil
}
