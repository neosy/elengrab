package dto

import (
	"github.com/google/uuid"
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
