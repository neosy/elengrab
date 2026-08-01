package dto

import (
	"time"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type TrackMediaWatchEventRequest struct {
	DownloadID uuid.UUID

	UserID    *uuid.UUID
	SessionID *uuid.UUID

	Position time.Duration
	Interval time.Duration

	EventType dtypes.MediaWatchEventType
}
