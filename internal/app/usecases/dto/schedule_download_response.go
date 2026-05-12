package dto

import (
	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/utils/hash"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type ScheduleDownloadResponse struct {
	URL string
	// file ID
	DownloadID uuid.UUID
	Status     dtypes.MediaDownloadStatus
	MediaTitle string
	// format type
	Format string
}

func (r *ScheduleDownloadResponse) ImageMetaHash(withValues ...any) string {
	values := []any{
		r.DownloadID,
		r.Status.String(),
	}

	if len(withValues) > 0 {
		values = append(values, withValues...)
	}

	return hash.MetaHashHex32(values)
}
