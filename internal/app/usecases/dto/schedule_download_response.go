package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/utils/hash"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type ScheduleDownloadResponse struct {
	URL string
	// file ID
	FileID     uuid.UUID
	Status     dtypes.FileStatus
	MediaTitle string
	// format type
	Format string
}

func (r *ScheduleDownloadResponse) ImageMetaHash(withValues ...any) string {
	values := []any{
		r.FileID,
		r.Status.String(),
		time.Now(),
	}

	if len(withValues) > 0 {
		values = append(values, withValues...)
	}

	return hash.MetaHashHex32(values)
}
