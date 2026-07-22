package linkweb

import (
	"context"

	dlink "github.com/neosy/elengrab/internal/domain/link"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

func (u *LinkWeb) FindLastByShortCode(ctx context.Context, shortCode string) (*dlink.Link, error) {
	return u.link.FindLastByShortCode(ctx, shortCode)
}

func (u *LinkWeb) GetLastByShortCode(ctx context.Context, shortCode string) (*dlink.Link, error) {
	return u.link.GetLastByShortCode(ctx, shortCode)
}

func (u *LinkWeb) GetLastShortURL(ctx context.Context, url string) (string, error) {
	shortCode := u.link.GenerateShortCodeByURL(url, u.options.ShortCodeLength, true)

	link, err := u.link.GetLastByShortCode(ctx, shortCode)
	if err != nil {
		return "", errorx.Errorf(
			"failed to get short link: %w", err,
			errorx.WithErrorMessage("Failed to get short link"),
		)
	}

	return link.ShortURL, nil
}
