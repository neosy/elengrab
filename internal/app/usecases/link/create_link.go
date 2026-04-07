package link

import (
	"context"
	"fmt"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	linkerr "github.com/neosy/elengrab/internal/app/usecases/link/errors"
	dlink "github.com/neosy/elengrab/internal/domain/link"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (u *Link) CreateLink(ctx context.Context, req *dto.LinkCreateRequest) (*dlink.Link, error) {
	if req == nil {
		return nil, linkerr.ErrFunctionNilParameter
	}

	// Checking fields from the request
	err := u.validateCreateRequestFields(req)
	if err != nil {
		u.logger.Error(
			"Ошибка при проверке полей из запроса",
			"error", err,
		)

		return nil, fmt.Errorf("error validating fields: %v", err)
	}

	// Transform the incoming request into a domain model
	link := u.mappers.MapLinkCreateRequestDtoToDomain(req)

	// Set the default base URL from the configuration
	baseURL := u.options.BaseURL
	// If base_url is explicitly specified in the request, we use it instead of the default value
	if req.BaseURL != nil && *req.BaseURL != "" {
		baseURL = *req.BaseURL
	}

	// Set the default shortcode length from the configuration
	shortCodeLength := u.options.ShortCodeLength
	// If the request specifies the short_code length and it is greater than zero, we use it
	if req.ShortCodeLength != nil && *req.ShortCodeLength > 0 {
		shortCodeLength = *req.ShortCodeLength
	}

	deterministic := u.options.Deterministic
	if req.Deterministic != nil {
		deterministic = *req.Deterministic
	}

	// Create a record
	link, err = u.link.Create(ctx, link, baseURL, shortCodeLength, deterministic)
	if err != nil {
		return nil, err
	}

	return link, err
}

func (u *Link) validateCreateRequestFields(req *dto.LinkCreateRequest) error {
	if req.OriginalURL == "" {
		u.logger.Error(
			"The originalURL field is empty.",
		)

		return errorx.New("originalURL is empty", exceptionx.VALIDATE)
	}

	return nil
}
