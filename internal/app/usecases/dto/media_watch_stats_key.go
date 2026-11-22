package dto

import "uuid"

type MediaWatchUpdateKey struct {
	DownloadID uuid.UUID
	UserID     uuid.UUID
	SessionID  uuid.UUID
}
