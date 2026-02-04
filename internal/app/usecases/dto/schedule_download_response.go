package dto

import (
	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type ScheduleDownloadResponse struct {
	// file id
	FileId     uuid.UUID
	Status     dtypes.FileStatus
	MediaTitle string
	// format type
	Format string
}
