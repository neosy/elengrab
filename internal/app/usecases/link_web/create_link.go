package linkweb

import (
	"context"
	"time"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

const (
	refreshThresholdDefault = 7 * 24 * time.Hour
)

func (u *linkWeb) CreateShortLink(ctx context.Context, url string) (string, error) {
	shortCode := u.link.GenerateShortCode(url, u.options.ShortCodeLength, true)

	link, err := u.link.FindLastByShortCode(ctx, shortCode)
	if err != nil {
		return "", errorx.Errorf(
			"failed to get short link: %w", err,
			errorx.WithErrorMessage("Failed to find short link"),
		)
	}

	shortURL := u.options.BaseShortURL + "/" + shortCode

	now := time.Now().UTC()

	if link != nil && link.ShortURL == shortURL {
		if link.ExpiresAt == nil {
			return link.ShortURL, nil
		}

		refreshThreshold := refreshThresholdDefault
		if u.options.RefreshThreshold != 0 {
			refreshThreshold = u.options.RefreshThreshold
		}

		// if link has more than shortCodeRefreshThreshold remaining, reuse existing short URL
		if !link.IsExpired(now.Add(refreshThreshold)) {
			return link.ShortURL, nil
		}
	}

	var expiresAt *time.Time
	if u.options.LinkTTL != 0 {
		expiresAt = new(now.Add(u.options.LinkTTL))
	}

	var shortCodeLength *uint8
	if u.options.ShortCodeLength != 0 {
		shortCodeLength = &u.options.ShortCodeLength
	}

	link, err = u.link.Create(
		ctx,
		&dto.LinkCreateRequest{
			OriginalURL:     url,
			ExpiresAt:       expiresAt,
			ShortCodeLength: shortCodeLength,
			Deterministic:   new(true),
		},
	)
	if err != nil {
		return "", errorx.Errorf(
			"failed to create short link: %w", err,
			errorx.WithErrorMessage("Failed to create short link"),
		)
	}

	return link.ShortURL, nil
}
