package ddownload

import (
	"github.com/google/uuid"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type RefreshMetadataTask struct {
	// Unique task identifier (UUID)
	TaskID uuid.UUID

	// Unique download media identifier (UUID)
	DownloadID uuid.UUID

	// UserID
	UserID uuid.UUID

	// Id worker
	WorkerID *uint64

	// ID job
	JobID *uuid.UUID
}

func (src *RefreshMetadataTask) Copy() *RefreshMetadataTask {
	if src == nil {
		return nil
	}

	copy := uptr.Copy(src)

	return copy
}
