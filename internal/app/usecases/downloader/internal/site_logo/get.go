package sitelogo

import (
	"context"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	"github.com/neosy/elengrab/pkg/errorx"
	"github.com/neosy/elengrab/pkg/errorx/exceptionx"
)

// FindByLogoID
// Record may not exist — caller decides what to do
func (uc *SiteLogo) FindByLogoID(ctx context.Context, logoID uuid.UUID) (*dmedia.SiteLogo, error) {
	if logoID == uuid.Nil {
		return nil, nil
	}

	logo, _ := uc.logoCacheRep.FindByLogoID(ctx, logoID)
	if logo != nil {
		return logo, nil
	}

	logo, err := uc.logoRep.FindByLogoID(ctx, logoID)
	if err != nil {
		uc.logger.Warn("Failed get siteLogo", "error", err)
		return nil, errorx.NewByErr(err, exceptionx.ERROR)
	}

	err = uc.logoCacheRep.Insert(ctx, logo)
	if err != nil {
		uc.logger.Warn("Failed to insert siteLogo cache", "error", err)
	}

	return logo, nil
}

// GetByLogoID
// Record MUST exist — otherwise NOT_FOUND
func (uc *SiteLogo) GetByLogoID(ctx context.Context, logoID uuid.UUID) (*dmedia.SiteLogo, error) {
	logo, err := uc.FindByLogoID(ctx, logoID)
	if err != nil {
		return nil, errorx.NewByErr(err, exceptionx.ERROR)
	}

	if logo == nil {
		uc.logger.Warn("Logo not found", "logoID", logoID)
		return nil, errorx.New("logo not found", exceptionx.NOT_FOUND)
	}

	return logo, nil
}

func (uc *SiteLogo) ExistsByLogoID(ctx context.Context, logoID uuid.UUID) (bool, error) {
	exists, _ := uc.logoCacheRep.ExistsByLogoID(ctx, logoID)
	if exists {
		return exists, nil
	}

	exists, err := uc.logoRep.ExistsByLogoID(ctx, logoID)
	if err != nil {
		uc.logger.Warn("Failed to check if site logo exists", "logoID", logoID, "error", err)
	}

	return exists, nil
}

// FindBySiteURL
// Record may not exist — caller decides what to do
func (uc *SiteLogo) FindBySiteURL(ctx context.Context, siteURL string) (*dmedia.SiteLogo, error) {
	if siteURL == "" {
		return nil, nil
	}

	logo, _ := uc.logoCacheRep.FindBySiteURL(ctx, siteURL)
	if logo != nil {
		return logo, nil
	}

	logo, err := uc.logoRep.FindBySiteURL(ctx, siteURL)
	if err != nil {
		uc.logger.Warn("Failed get siteLogo", "error", err)
		return nil, errorx.NewByErr(err, exceptionx.ERROR)
	}

	if logo == nil {
		return nil, nil
	}

	err = uc.logoCacheRep.Insert(ctx, logo)
	if err != nil {
		uc.logger.Warn("Failed to insert siteLogo cache", "error", err)
	}

	return logo, nil
}

// GetBySiteURL
// Record MUST exist — otherwise NOT_FOUND
func (uc *SiteLogo) GetBySiteURL(ctx context.Context, siteURL string) (*dmedia.SiteLogo, error) {
	logo, err := uc.FindBySiteURL(ctx, siteURL)
	if err != nil {
		return nil, errorx.NewByErr(err, exceptionx.ERROR)
	}

	if logo == nil {
		uc.logger.Warn("Logo not found", "siteURL", siteURL)
		return nil, errorx.New("logo not found", exceptionx.NOT_FOUND)
	}

	return logo, nil
}

func (uc *SiteLogo) ExistsBySiteURL(ctx context.Context, siteURL string) (bool, error) {
	exists, _ := uc.logoCacheRep.ExistsBySiteURL(ctx, siteURL)
	if exists {
		return exists, nil
	}

	exists, err := uc.logoRep.ExistsBySiteURL(ctx, siteURL)
	if err != nil {
		uc.logger.Warn("Failed to check if site logo exists", "siteURL", siteURL, "error", err)
	}

	return exists, nil
}

func (uc *SiteLogo) ExistsBySiteURLFromCache(ctx context.Context, siteURL string) (bool, error) {
	if siteURL == "" {
		return false, nil
	}

	return uc.logoCacheRep.ExistsBySiteURL(ctx, siteURL)
}
