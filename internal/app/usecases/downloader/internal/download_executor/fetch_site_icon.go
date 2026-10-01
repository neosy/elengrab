package dlexecutor

import (
	"context"
	"time"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
)

func (uc *Executor) fetchAndUpdateSiteLogo(ctx context.Context, url string) error {
	baseURL := httpx.BaseURL(url)

	logo, err := uc.siteIcon.FindBySiteURLWithoutCache(ctx, baseURL)
	if err != nil {
		return err
	}

	if logo != nil && time.Since(logo.UpdatedAt) <= uc.logoUpdateInterval {
		return nil
	}

	fetchedLogo, err := uc.fetchSiteLog(ctx, baseURL)
	if err != nil {
		return err
	}

	if fetchedLogo == nil {
		return nil
	}

	if logo != nil {
		if logo.Equal(fetchedLogo) {
			return nil
		}

		logo.SetRequired(baseURL, fetchedLogo.SiteTitle, fetchedLogo.ImageData())

		return uc.siteIcon.Update(ctx, logo)
	}

	logo = dmedia.NewSiteLogo(baseURL, fetchedLogo.SiteTitle, fetchedLogo.ImageData())

	return uc.siteIcon.Create(ctx, logo)
}

func (uc *Executor) fetchSiteLog(ctx context.Context, url string) (*dmedia.SiteLogo, error) {
	// Fetch the site title.
	startTime := time.Now()
	title, err := httpx.FetchTitle(
		ctx,
		url,
		httpx.ClientOptionWithTimeout(getHTMLTimeout),
		httpx.ClientOptionWithDefaultCookieJar(),
	)
	elapsed := time.Since(startTime)
	if err != nil {
		uc.logger.Debug(
			"Failed to fetch site title",
			"url", url,
			"elapsed", uformat.DurationFormat(elapsed),
			"error", err)
	}
	if title != "" {
		uc.logger.Debug(
			"Site title fetched",
			"title", title,
			"url", url,
			"elapsed", uformat.DurationFormat(elapsed),
		)
	}

	// Fetch the best icon.
	startTime = time.Now()
	image, err := uc.siteIconFetcher.FetchBestIcon(ctx, url)
	elapsed = time.Since(startTime)
	if err != nil {
		uc.logger.Warn(
			"Failed to fetch icon",
			"url", url,
			"elapsed", uformat.DurationFormat(elapsed),
			"error", err,
		)
		return nil, err
	}
	uc.logger.Debug(
		"Best icon fetched",
		"url", url,
		"elapsed", uformat.DurationFormat(elapsed),
	)

	// Create new siteLogo
	logo := dmedia.NewSiteLogo(url, title, image)

	return logo, nil
}
