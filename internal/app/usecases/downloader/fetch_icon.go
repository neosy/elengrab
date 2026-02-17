package downloader

import (
	"context"
	"time"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	"github.com/neosy/elengrab/pkg/httpx"
)

func (uc *YouTubeDownloader) fetchIcon(ctx context.Context, url string) error {
	baseURL := httpx.BaseURL(url)

	logo, err := uc.siteIcon.FindBySiteURL(ctx, baseURL)
	if err != nil {
		return err
	}

	if logo != nil && time.Since(logo.UpdatedAt) <= uc.logoUpdateInterval {
		return nil
	}

	// Fetch the site title.
	startTime := time.Now()
	title, err := httpx.GetTitle(
		ctx,
		baseURL,
		httpx.ClientOptionWithTimeout(getHTMLTimeout),
		httpx.ClientOptionWithDefaultCookieJar(),
	)
	elapsed := time.Since(startTime)
	if err != nil {
		uc.logger.Debug(
			"Failed to fetch site title",
			"baseURL", baseURL,
			"elapsed", elapsed,
			"error", err)
	}
	if title != "" {
		uc.logger.Debug(
			"Site title fetched",
			"title", title,
			"baseURL", baseURL,
			"elapsed", elapsed,
		)
	}

	// Fetch the best icon.
	startTime = time.Now()
	image, err := uc.siteIconFetcher.FetchBestIcon(ctx, baseURL)
	elapsed = time.Since(startTime)
	if err != nil {
		uc.logger.Warn(
			"Failed to fetch icon",
			"url", baseURL,
			"elapsed", elapsed,
			"error", err,
		)
		return err
	}
	uc.logger.Debug(
		"Best icon fetched",
		"url", baseURL,
		"elapsed", elapsed,
	)

	// Update existing logo if it's outdated.
	if logo != nil {
		logo.SetRequired(baseURL, title, image)

		err := uc.siteIcon.Update(ctx, logo)
		if err != nil {
			return err
		}
		return nil
	}

	// Create new siteLogo
	logo = dmedia.NewSiteLogo(baseURL, title, image)

	err = uc.siteIcon.Create(ctx, logo)
	if err != nil {
		return err
	}

	return nil
}
