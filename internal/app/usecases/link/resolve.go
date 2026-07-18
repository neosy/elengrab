package link

import (
	"context"
	"time"

	dlink "github.com/neosy/elengrab/internal/domain/link"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/valyala/fasthttp"
)

func (u *Link) ResolveShortCode(ctx context.Context, shortCode string) (*dlink.Link, error) {
	link, err := u.link.FindLastByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	if link == nil {
		return nil, errorx.NewHTTPMessage("link not found", fasthttp.StatusGone)
	}

	now := time.Now().UTC()

	// Check if the link has expired
	if link.IsExpired(now) {
		return nil, errorx.NewHTTPMessage("the link has expired and is no longer valid", fasthttp.StatusGone)
	}

	return link, nil
}
