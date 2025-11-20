package dto

import (
	"time"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type GetFileInfoResponse struct {
	FileId               uuid.UUID
	Status               dtypes.FileStatus
	YoutubeUrl           string
	YoutubeTitle         string
	FileName             string
	FileExt              string
	FileFullName         string
	FilePath             string
	FileSize             *int
	SafeReadableFullName string
	StatusText           string
	CreatedAt            time.Time
}
