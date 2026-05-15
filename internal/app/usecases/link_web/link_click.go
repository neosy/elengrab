package linkweb

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dlink "github.com/neosy/elengrab/internal/domain/link"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

func (u *LinkWeb) ShortLinkClick(
	ctx context.Context,
	shortURL string,
	ipAddress string,
	userAgent string,
	referrer string,
) (*dlink.Link, error) {
	var userAgentPtr, referrerPtr *string

	if userAgent != "" {
		userAgentPtr = &userAgent
	}

	if referrer != "" {
		referrerPtr = &referrer
	}

	link, err := u.link.Click(
		ctx,
		&dto.ShortLinkClickRequest{
			ShortURL:  shortURL,
			IPAddress: ipAddress,
			UserAgent: userAgentPtr,
			Referrer:  referrerPtr,
		},
	)
	if err != nil {
		return nil, errorx.Errorf("failed to process click on short link: %w", err)
	}

	return link, nil
}
