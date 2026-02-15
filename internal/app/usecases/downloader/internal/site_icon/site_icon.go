package siteicon

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type SiteIcon struct {
	logger *slog.Logger

	// repositories
	logoRep persistence.SiteLogoRepository

	// cache
	logoCacheRep persistence.SiteLogoCacheRepository
}

func NewSiteIcon(
	logger *slog.Logger,
	logoRep persistence.SiteLogoRepository,
	logoCacheRep persistence.SiteLogoCacheRepository,
) *SiteIcon {
	return &SiteIcon{
		logger:       logger,
		logoRep:      logoRep,
		logoCacheRep: logoCacheRep,
	}
}
