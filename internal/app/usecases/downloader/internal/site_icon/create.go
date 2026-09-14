package siteicon

import (
	"context"
	"errors"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

// Create creates a new site logo entry.
func (uc *SiteIcon) Create(ctx context.Context, logo *dmedia.SiteLogo) error {
	if logo == nil {
		uc.logger.Warn("Nil pointer in function")
		return errors.New("function parameter is a nil pointer")
	}

	// Ensure LogoID is generated if it's not provided
	if logo.LogoID == uuid.Nil {
		logo.LogoID = uuid.New()
	}

	defer uc.logoCacheRep.Delete(ctx, logo.LogoID)

	// Attempt to insert the site logo into the repository
	err := uc.logoRepo().Insert(ctx, logo)
	if err != nil {
		uc.logger.Error(
			"Failed to insert record into siteLogo cache repository",
			"logoURL", logo.ImageURL,
			"error", err,
		)
		return err
	}

	return nil
}
