package logofetcher

import "log/slog"

type SiteLogoFetcher struct {
	logger *slog.Logger
}

func NewSiteLogoFetcher(logger *slog.Logger) *SiteLogoFetcher {
	return &SiteLogoFetcher{
		logger: logger,
	}
}
