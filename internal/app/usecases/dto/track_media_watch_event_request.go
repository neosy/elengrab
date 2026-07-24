package dto

import (
	"time"

	"github.com/google/uuid"
)

type TrackMediaWatchEventRequest struct {
	DownloadID uuid.UUID

	UserID    *uuid.UUID
	SessionID *uuid.UUID

	Position time.Duration
	Interval time.Duration
}
