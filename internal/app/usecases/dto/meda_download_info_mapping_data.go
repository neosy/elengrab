package dto

import (
	"time"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
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

	Channel *dmedia.Channel

	Progress *dservices.DownloaderProgress
}
