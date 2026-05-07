package iconfetcher

import (
	"context"
	"errors"
	"net/url"
	"sort"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/httpx"
)

var ErrIconNotFound = errors.New("icon not found")

// FetchBestIcon fetches the best available icon for the given site URL.
// It retrieves icon candidates from the site, downloads them in priority order,
// and returns the first successfully downloaded image.
// Returns ErrIconNotFound if no valid icon could be retrieved.
func (lf *SiteIconFetcher) FetchBestIcon(ctx context.Context, siteURL string) (*dtypes.ImageData, error) {
	candidates, err := lf.fetchCandidatesFromURL(ctx, siteURL)
	if err != nil {
		return nil, err
	}

	for _, c := range candidates {
		image, err := lf.fetchImage(ctx, c.url)
		if err == nil {
			return image, nil
		}
	}

	return nil, ErrIconNotFound
}

// FetchIcons fetches all available icons for the given site URL.
// It retrieves icon candidates from the site, downloads each candidate,
// and returns all successfully downloaded images.
func (lf *SiteIconFetcher) FetchIcons(ctx context.Context, siteURL string) ([]*dtypes.ImageData, error) {
	candidates, err := lf.fetchCandidatesFromURL(ctx, siteURL)
	if err != nil {
		return nil, err
	}

	capHint := min(len(candidates), 4)
	icons := make([]*dtypes.ImageData, 0, capHint)
	for _, c := range candidates {
		icon, err := lf.fetchImage(ctx, c.url)
		if err == nil {
			icons = append(icons, icon)
		}
	}

	return icons, nil
}

// fetchCandidatesFromURL retrieves icon candidates from the site's HTML head,
// applies a fallback to /favicon.ico, and sorts candidates by priority.
func (lf *SiteIconFetcher) fetchCandidatesFromURL(ctx context.Context, siteURL string) ([]iconCandidate, error) {
	links, err := httpx.GetLinksInHead(
		ctx,
		siteURL,
		httpx.MethodGetOptions{
			Limit:            limitHTML,
			IgnoreStatusCode: true,
		},
		httpx.ClientOptionWithTimeout(fetchIconTimeout),
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
