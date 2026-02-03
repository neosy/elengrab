package downloader

import (
	"context"
	"fmt"
	"time"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	"github.com/neosy/elengrab/pkg/httpx"
)

func (uc *YouTubeDownloader) fetchLogo(ctx context.Context, url string) error {
	baseURL := httpx.GetBaseURL(url)
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
		if time.Since(logo.UpdatedAt) <= logoUpdateInterval {
			return nil
		}

		err := uc.siteLogo.Update(ctx, logo)
		if err != nil {
			return err
		}
		return nil
	}

	logo = &dmedia.SiteLogo{
		SiteURL:     baseURL,
		ImageURL:    image.URL,
		ImageRaw:    image.Raw,
		ImageFormat: image.Format,
	}

	err = uc.siteLogo.Create(ctx, logo)
	if err != nil {
		return err
	}

	return nil
}
