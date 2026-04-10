package link

import (
	"context"
	"slices"
	"time"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dlink "github.com/neosy/elengrab/internal/domain/link"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
	"github.com/valyala/fasthttp"
)

// ShortLinkClick handles a click on a short link: validates the request, checks link availability,
// and records the click in the database. Returns the created click record.
func (u *Link) ShortLinkClick(ctx context.Context, req *dto.ShortLinkClickRequest) (*dlink.Link, error) {
	linkClick := u.mappers.MapShortLinkClickRequestToDomain(req)
	return u.click(ctx, linkClick)
}

func (u *Link) click(
	ctx context.Context,
	linkClickDraft *dlink.LinkClick,
) (*dlink.Link, error) {
	// Copy the draft to avoid modifying the original
	linkClick := uptr.Copy(linkClickDraft)

	// Extract shortCode from the provided URL
	shortCode := dlink.GetShortCodeFromURL(linkClick.ShortURL)
	if shortCode == "" {
		return nil, errorx.New(
			"shortCode is not correct",
			exceptionx.VALIDATE,
			errorx.WithErrorMessage("Invalid short code"),
		)
	}

	// Find the latest active link by shortCode
	link, err := u.link.GetLastByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	linkClick.LinkID = link.LinkID

	// Validate if the link can be clicked
	err = u.validateClick(ctx, link, linkClick)
	if err != nil {
		return nil, err
	}

	// Create a record of the link click
	_, err = u.linkClick.Create(ctx, linkClick)
	if err != nil {
		return nil, err
	}

	return link, err
}

func (u *Link) validateClick(
	ctx context.Context,
	link *dlink.Link,
	linkClick *dlink.LinkClick,
) error {
	// Check if the link has been soft-deleted (DeletedAt set)
	if link.DeletedAt != nil {
		return errorx.NewHTTPMessage("this link has expired", fasthttp.StatusGone)
	}

	// Check if the link has expired
	if link.ExpiresAt != nil && link.ExpiresAt.Before(time.Now()) {
		return errorx.NewHTTPMessage("the link has expired and is no longer valid", fasthttp.StatusGone)
	}

	// Check full match of ShortURL if required
	if link.IsMatchShortURL && link.ShortURL != linkClick.ShortURL {
		return errorx.NewHTTP(
			"short URL does not match the expected value",
			fasthttp.StatusBadRequest,
			errorx.WithErrorMessage("Invalid short URL"),
		)
	}

	// Check click limit
	if link.MaxClicks != nil && *link.MaxClicks > 0 {
		count, err := u.linkClick.CountByLinkId(ctx, link.LinkID)
		if err != nil {
			return err
		}

		// Return error if maximum clicks reached
		if count >= *link.MaxClicks {
			return errorx.NewHTTP(
				"maximum number of clicks reached",
				fasthttp.StatusGone,
				errorx.WithErrorMessage("This link has reached its usage limit"),
			)
		}
	}

	// Check if the user's IP address is allowed
	if len(link.AllowedIPs) > 0 {
		var exists bool
		for _, ip := range link.AllowedIPs {
			if ip == linkClick.IPAddress {
				exists = true
				break
			}
		}
		if !exists {
			return errorx.NewHTTP(
				"access from this IP address is not allowed",
				fasthttp.StatusForbidden,
				errorx.WithErrorMessage("Access denied from this IP address"),
			)
		}
	}

	// Check if the user is allowed to access the link
	if linkClick.ClickedBy != nil && len(link.AllowedUserIDs) > 0 {
		if !slices.Contains(link.AllowedUserIDs, *linkClick.ClickedBy) {
			return errorx.NewHTTP(
				"this user is not allowed to access the link",
				fasthttp.StatusForbidden,
				errorx.WithErrorMessage("You are not allowed to access this link"),
			)
		}
	}

	return nil
}
