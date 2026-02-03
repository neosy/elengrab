package sitelogo

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type SiteLogo struct {
	logger *slog.Logger

	// repositories
	logoRep persistence.SiteLogoRepository

	// cache
	logoCacheRep persistence.SiteLogoCacheRepository
}

func NewSiteLogo(
	logger *slog.Logger,
	logoRep persistence.SiteLogoRepository,
	logoCacheRep persistence.SiteLogoCacheRepository,
) *SiteLogo {
	return &SiteLogo{
		logger:       logger,
		logoRep:      logoRep,
		logoCacheRep: logoCacheRep,
	}
}
