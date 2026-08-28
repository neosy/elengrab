package siteicon

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type SiteIcon struct {
	logger *slog.Logger

	// repositories
	logoRepo persistence.SiteLogoRepositoryFactory

	// cache
	logoCacheRep persistence.SiteLogoCacheRepository
}

func NewSiteIcon(
	logger *slog.Logger,
	logoRepo persistence.SiteLogoRepositoryFactory,
	logoCacheRep persistence.SiteLogoCacheRepository,
) *SiteIcon {
	return &SiteIcon{
		logger:       logger,
		logoRepo:     logoRepo,
		logoCacheRep: logoCacheRep,
	}
}
