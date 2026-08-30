package link

import (
	"context"

	dlink "github.com/neosy/elengrab/internal/domain/link"
)

func (u *link) FindLastByShortCode(ctx context.Context, shortCode string) (*dlink.Link, error) {
	return u.link.FindLastByShortCode(ctx, shortCode)
}

func (u *link) GetLastByShortCode(ctx context.Context, shortCode string) (*dlink.Link, error) {
	return u.link.GetLastByShortCode(ctx, shortCode)
}
