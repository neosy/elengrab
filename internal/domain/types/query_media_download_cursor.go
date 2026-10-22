package dtypes

import (
	"time"

	"uuid"
)

type QueryMediaDownloadCursor struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Views     uint32
}
