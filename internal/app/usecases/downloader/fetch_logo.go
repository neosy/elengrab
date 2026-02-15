package downloader

import (
	"context"
	"time"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	"github.com/neosy/elengrab/pkg/httpx"
)

func (uc *YouTubeDownloader) fetchLogo(ctx context.Context, url string) error {
	title, err := httpx.GetTitle(
		ctx,
		httpx.BaseURL(url),
		httpx.ClientOptionWithTimeout(getHTMLTimeout),
	)
	if err != nil {
		uc.logger.Warn("Failed to get title of site", "error", err)
	}

	baseURL := httpx.BaseURL(url)

	logo, err := uc.siteLogo.FindBySiteURL(ctx, baseURL)
	if err != nil {
		return err
	}

	if logo != nil && time.Since(logo.UpdatedAt) <= uc.logoUpdateInterval {
		return nil
	}

	image, err := uc.siteLogoFetcher.FetchBestLogo(ctx, baseURL)
	if err != nil {
		uc.logger.Warn("Failed to fetch logo", "url", baseURL, "error", err)
		return err
	}

	// Update existing logo if it's outdated.
	if logo != nil {
		logo.SetRequired(baseURL, title, image)

		err := uc.siteLogo.Update(ctx, logo)
		if err != nil {
			return err
		}
		return nil
	}

	// Create new siteLogo
	logo = dmedia.NewSiteLogo(baseURL, title, logo.ImageData())

	err = uc.siteLogo.Create(ctx, logo)
	if err != nil {
		return err
	}

	return nil
}
