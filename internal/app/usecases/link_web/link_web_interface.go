package linkweb

import (
	"context"

	dlink "github.com/neosy/elengrab/internal/domain/link"
)

type LinkWeb interface {
	CreateShortLink(ctx context.Context, url string) (string, error)
	DeleteShortLink(ctx context.Context, url string) error

	FindLastShortURL(ctx context.Context, url string) (string, error)
	GetLastShortURL(ctx context.Context, url string) (string, error)

	GetLastByShortCode(ctx context.Context, shortCode string) (*dlink.Link, error)

	ResolveURL(ctx context.Context, url string) (*dlink.Link, error)

	ShortLinkClick(
		ctx context.Context,
		shortURL string,
		ipAddress string,
		userAgent string,
		referrer string,
	) (*dlink.Link, error)
}
