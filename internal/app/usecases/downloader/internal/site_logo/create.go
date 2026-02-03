package sitelogo

import (
	"context"
	"errors"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

func (uc *SiteLogo) Create(ctx context.Context, logo *dmedia.SiteLogo) error {
	if logo == nil {
		uc.logger.Warn("Nil pointer in function")
		return errors.New("function parameter is a nil pointer")
	}

	if logo.LogoID == uuid.Nil {
		logo.LogoID = uuid.New()
	}

	err := uc.logoRep.Insert(ctx, logo)
	if err != nil {
		uc.logger.Warn(
			"Failed to insert record into siteLogo cache repository",
			"error", err,
		)
		return err
	}

	return nil
}
