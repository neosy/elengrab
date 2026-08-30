package linkweb

import (
	"context"

	dlink "github.com/neosy/elengrab/internal/domain/link"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

func (u *linkWeb) FindLastByShortCode(ctx context.Context, shortCode string) (*dlink.Link, error) {
	return u.link.FindLastByShortCode(ctx, shortCode)
}

func (u *linkWeb) GetLastByShortCode(ctx context.Context, shortCode string) (*dlink.Link, error) {
	return u.link.GetLastByShortCode(ctx, shortCode)
}

func (u *linkWeb) FindLastShortURL(ctx context.Context, url string) (string, error) {
	shortCode := u.link.GenerateShortCode(url, u.options.ShortCodeLength, true)

	link, err := u.link.FindLastByShortCode(ctx, shortCode)
	if err != nil {
		return "", errorx.Errorf(
			"failed to get short link: %w", err,
			errorx.WithErrorMessage("Failed to get short link"),
		)
	}

	if link == nil {
		return "", nil
	}

	return link.ShortURL, nil
}

func (u *linkWeb) GetLastShortURL(ctx context.Context, url string) (string, error) {
	shortCode := u.link.GenerateShortCode(url, u.options.ShortCodeLength, true)

	link, err := u.link.GetLastByShortCode(ctx, shortCode)
	if err != nil {
		return "", errorx.Errorf(
			"failed to get short link: %w", err,
			errorx.WithErrorMessage("Failed to get short link"),
		)
	}

	return link.ShortURL, nil
}
