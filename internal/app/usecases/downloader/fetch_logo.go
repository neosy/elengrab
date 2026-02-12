package downloader

import (
	"context"
	"fmt"
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
	image, err := uc.siteLogoFetcher.FetchLogo(ctx, baseURL)
	if err != nil {
		uc.logger.Warn("Failed to fetch logo", "url", baseURL, "error", err)
		return err
	}
	if image == nil {
		uc.logger.Warn("Site log not found", "url", baseURL)
		return fmt.Errorf("site logo not found")
	}

	logo, _ := uc.siteLogo.FindBySiteURL(ctx, baseURL)
	if logo != nil {
		if time.Since(logo.UpdatedAt) <= uc.logoUpdateInterval {
			return nil
		}

		// Update existing logo if it's outdated.
		logo.SetRequired(baseURL, title, image)

		err := uc.siteLogo.Update(ctx, logo)
		if err != nil {
			return err
		}
		return nil
	}

	// Create new siteLogo
	logo = dmedia.NewSiteLogo(baseURL, title, image)

	err = uc.siteLogo.Create(ctx, logo)
	if err != nil {
		return err
	}

	return nil
}
