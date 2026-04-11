package siteicon

import (
	"context"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

// FindByLogoID
// Record may not exist — caller decides what to do
func (uc *SiteIcon) FindByLogoID(ctx context.Context, logoID uuid.UUID) (*dmedia.SiteLogo, error) {
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
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	err = uc.logoCacheRep.Save(ctx, logo)
	if err != nil {
		uc.logger.Warn("Failed to insert siteLogo cache", "error", err)
	}

	return logo, nil
}

// GetByLogoID
// Record MUST exist — otherwise NOT_FOUND
func (uc *SiteIcon) GetByLogoID(ctx context.Context, logoID uuid.UUID) (*dmedia.SiteLogo, error) {
	logo, err := uc.FindByLogoID(ctx, logoID)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if logo == nil {
		uc.logger.Warn("Logo not found", "logoID", logoID)
		return nil, errorx.New("logo not found", exceptionx.NOT_FOUND)
	}

	return logo, nil
}

// Checks if a site logo exists by its unique identifier.
func (uc *SiteIcon) ExistsByLogoID(ctx context.Context, logoID uuid.UUID) (bool, error) {
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
func (uc *SiteIcon) FindBySiteURL(ctx context.Context, siteURL string) (*dmedia.SiteLogo, error) {
	if siteURL == "" {
		return nil, nil
	}

	logo, cacheStatus, _ := uc.logoCacheRep.FindBySiteURL(ctx, siteURL)
	if logo != nil {
		return logo, nil
	}
	if cacheStatus == memsimple.CacheStatusNegativeHit {
		return nil, nil
	}

	logo, err := uc.logoRep.FindBySiteURL(ctx, siteURL)
	if err != nil {
		uc.logger.Warn("Failed get siteLogo", "error", err)
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if logo == nil {
		err = uc.logoCacheRep.SaveNegative(ctx, siteURL)
		return nil, nil
	}

	err = uc.logoCacheRep.Save(ctx, logo)
	if err != nil {
		uc.logger.Warn("Failed to insert siteLogo cache", "error", err)
	}

	return logo, nil
}

// GetBySiteURL
// Record MUST exist — otherwise NOT_FOUND
func (uc *SiteIcon) GetBySiteURL(ctx context.Context, siteURL string) (*dmedia.SiteLogo, error) {
	logo, err := uc.FindBySiteURL(ctx, siteURL)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if logo == nil {
		return nil, errorx.New("logo not found", exceptionx.NOT_FOUND)
	}

	return logo, nil
}

// ExistsBySiteURL returns true if siteURL exists in the repository.
func (uc *SiteIcon) ExistsBySiteURL(ctx context.Context, siteURL string) (bool, error) {
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

// ExistsBySiteURLFromCache returns true if siteURL exists in the cache.
func (uc *SiteIcon) ExistsBySiteURLFromCache(ctx context.Context, siteURL string) (bool, error) {
	if siteURL == "" {
		return false, nil
	}

	return uc.logoCacheRep.ExistsBySiteURL(ctx, siteURL)
}
