package sitelogo

import (
	"context"
	"errors"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

// Create creates a new site logo entry.
func (uc *SiteLogo) Create(ctx context.Context, logo *dmedia.SiteLogo) error {
	if logo == nil {
		uc.logger.Warn("Nil pointer in function")
		return errors.New("function parameter is a nil pointer")
	}

	// Ensure LogoID is generated if it's not provided
	if logo.LogoID == uuid.Nil {
		logo.LogoID = uuid.New()
	}

	// Attempt to insert the site logo into the repository
	err := uc.logoRep.Insert(ctx, logo)
	if err != nil {
		uc.logger.Error(
			"Failed to insert record into siteLogo cache repository",
			"logoURL", logo.ImageURL,
			"error", err,
		)
		return err
	}

	logo, _ = uc.logoRep.FindByLogoID(ctx, logo.LogoID)
	if logo != nil {
		err := uc.logoCacheRep.Save(ctx, logo)
		if err != nil {
			uc.logger.Error(
				"Failed to save siteLogo cache",
				"logoURL", logo.ImageURL,
				"error", err,
			)
			return err
		}
	}

	return nil
}
