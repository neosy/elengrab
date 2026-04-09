package linkweb

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dlink "github.com/neosy/elengrab/internal/domain/link"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
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

	link, err := u.link.ShortLinkClick(
		ctx,
		&dto.ShortLinkClickRequest{
			ShortURL:  shortURL,
			IPAddress: ipAddress,
			UserAgent: userAgentPtr,
			Referrer:  referrerPtr,
		},
	)
	if err != nil {
		message := "Failed to process click on short link"
		if errorx.IsErrorx(err) {
			exception := errorx.OuterException(err)
			if exception != nil && exception.Num() == uint(exceptionx.NOT_FOUND) {
				message = "Short link not found."
			}
		}
		return nil, errorx.Errorf(
			"failed to process click on short link: %w", err,
			errorx.ErrorMessageArg(message),
		)
	}

	return link, nil
}
