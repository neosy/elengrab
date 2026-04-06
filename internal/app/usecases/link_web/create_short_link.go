package linkweb

import (
	"context"
	"time"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

const (
	expirationDays            = 14
	expirationDuration        = expirationDays * 24 * time.Hour
	shortCodeRefreshThreshold = 7 * 24 * time.Hour
)

func (u *LinkWeb) CreateShortLink(ctx context.Context, url string) (string, error) {
	shortCode := u.link.GenerateShortCodeByURL(url)

	link, err := u.link.GetLastByShortCode(ctx, shortCode)
	if err != nil {
		return "", errorx.Errorf(
			"failed to get short link: %w", err,
			errorx.ErrorMessageArg("Failed to find short link"),
		)
	}

	// if link has more than shortCodeRefreshThreshold remaining, reuse existing short URL
	if link != nil && link.ExpiresAt != nil && link.ExpiresAt.After(time.Now().UTC().Add(shortCodeRefreshThreshold)) {
		return link.ShortURL, nil
	}

	link, err = u.link.CreateLink(
		ctx,
		&dto.LinkCreateRequest{
			OriginalURL:     url,
			ExpiresAt:       uptr.Any(time.Now().Add(expirationDuration).UTC()),
			ShortCodeLength: &u.shortCodeLength,
		},
	)
	if err != nil {
		return "", errorx.Errorf(
			"failed to create short link: %w", err,
			errorx.ErrorMessageArg("Failed to create short link"),
		)
	}

	return link.ShortURL, nil
}
