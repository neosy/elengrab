package dto

import (
	"time"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type CreateMediaWatchEventRequest struct {
	Event *ddownload.MediaWatchEvent

	EventType     dtypes.MediaWatchEventType
	MediaDuration time.Duration
}
