package siteicon

import (
	"context"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

func (uc *SiteIcon) Update(ctx context.Context, logo *dmedia.SiteLogo) error {
	err := uc.logoRep.Update(ctx, logo)
	if err != nil {
		uc.logger.Warn("Update record error", "error", err)
		return err
	}

	if err := uc.logoCacheRep.Save(ctx, logo); err != nil {
		uc.logger.Warn("Update siteLogo cache error", "error", err)
		return err
	}

	return err
}
