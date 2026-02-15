package iconfetcher

import "log/slog"

type SiteIconFetcher struct {
	logger *slog.Logger
}

func NewSiteIconFetcher(logger *slog.Logger) *SiteIconFetcher {
	return &SiteIconFetcher{
		logger: logger,
	}
}
