package dto

import (
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
	SafeReadableFullName string
	StatusText           string
}
