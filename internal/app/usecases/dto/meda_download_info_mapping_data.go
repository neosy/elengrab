package dto

import (
	"time"

	dservices "github.com/neosy/elengrab/internal/domain/services"
)

type MediaDownloadInfoMappingData struct {
	UserLogin   string
	AvatarTitle string

	ViewCount uint32

	UserLastWatchPosition time.Duration
	UserWatched           bool

	HasSiteIcon         bool
	ThumbnailIsPortrait bool

	Progress *dservices.DownloaderProgress
}
