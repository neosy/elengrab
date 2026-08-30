package linkweb

import (
	"context"

	dlink "github.com/neosy/elengrab/internal/domain/link"
)

func (u *linkWeb) ResolveURL(ctx context.Context, url string) (*dlink.Link, error) {
	shortCode := u.link.GenerateShortCode(url, u.options.ShortCodeLength, true)
	return u.link.ResolveShortCode(ctx, shortCode)
}
