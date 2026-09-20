package dto

import (
	"time"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type MediaDownloadInfoMappingData struct {
	UserType  dtypes.UserType
	UserLogin string

	ViewCount uint32

	UserLastWatchPosition time.Duration
	UserWatched           bool

	HasSiteIcon         bool
	ThumbnailIsPortrait bool

	Channel      *dmedia.Channel
	ChannelTitle string

	Progress *dservices.DownloaderProgress
}
