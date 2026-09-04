package dtypes

import (
	"time"

	"github.com/google/uuid"
)

type QueryMediaDownloadCursor struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Views     uint32
}
