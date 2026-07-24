package dto

import (
	"time"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type CreateMediaWatchEventRequest struct {
	Event *ddownload.MediaWatchEvent

	MediaDuration time.Duration
}
