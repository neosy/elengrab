package link

import (
	"context"
	"strings"

	"github.com/google/uuid"
	linkerr "github.com/neosy/elengrab/internal/app/usecases/link/errors"
	dlink "github.com/neosy/elengrab/internal/domain/link"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

// maxGenShortCodeAttempts defines the maximum number of attempts
const maxGenShortCodeAttempts = 2

func (u *Link) Create(
	ctx context.Context,
	link *dlink.Link,
	baseURL string,
	shortCodeLength uint8,
	deterministic bool,
) (*dlink.Link, error) {
	if link == nil {
		u.logger.Warn("Nil pointer passed to function")
		return nil, linkerr.ErrFunctionNilParameter
	}

	// Generate a new UUID if not already set
	if link.LinkID == uuid.Nil {
		link.LinkID = uuid.New()
	}

	var existsShortCode bool

	if !deterministic {
		// Attempt to generate a unique short code (shortCode)
		// No more than maxGenShortCodeAttempts times
		for range maxGenShortCodeAttempts {
			// Generate the short code based on link ID and original URL
			link.ShortCode = dlink.GenerateShortCode(link.LinkID, link.OriginalURL, shortCodeLength, deterministic)

			// Check if this short code already exists
			existsShortCode, _ = u.linkRep.ExistsActiveShortCode(ctx, link.ShortCode)
			if !existsShortCode {
				break
			}
		}

		// If a unique shortCode could not be generated after all attempts
		if existsShortCode {
			u.logger.Warn(
				"Failed to generate a unique shortCode",
				"originalURL", link.OriginalURL,
				"attempts", maxGenShortCodeAttempts,
			)
		}
	} else {
		link.ShortCode = dlink.GenerateShortCode(link.LinkID, link.OriginalURL, shortCodeLength, deterministic)
	}

	// Build the full short URL
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	link.ShortURL = baseURL + link.ShortCode

	// Insert the link into the repository
	err := u.linkRep.Insert(ctx, link)
	if err != nil {
		u.logger.Warn(
			"Failed to insert record into repository",
			"error", err,
		)
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	// Retrieve the link from the repository to return
	link, err = u.GetByLinkID(ctx, link.LinkID)
	if err != nil {
		u.logger.Warn(
			"Failed to get record from repository",
			"error", err,
		)
		return nil, err
	}

	return link, nil
}
