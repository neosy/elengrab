package dto

import (
	"time"
	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type CreateMediaWatchEventRequest struct {
	Event *ddownload.MediaWatchEvent

	EventType     dtypes.MediaWatchEventType
	MediaDuration time.Duration
}

func (req *CreateMediaWatchEventRequest) MediaWatchUpdateKey() MediaWatchUpdateKey {
	if req == nil || req.Event == nil {
		return MediaWatchUpdateKey{}
	}

	var userID uuid.UUID
	if req.Event.UserID != nil {
		userID = *req.Event.UserID
	}

	var sessionID uuid.UUID
	if req.Event.SessionID != nil {
		sessionID = *req.Event.SessionID
	}

	return MediaWatchUpdateKey{
		DownloadID: req.Event.DownloadID,
		UserID:     userID,
		SessionID:  sessionID,
	}
}
