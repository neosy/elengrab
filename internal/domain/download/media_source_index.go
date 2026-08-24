package ddownload

import (
	"time"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type MediaSourceIndex struct {
	// Identifier of the watched media (UUID)
	DownloadID uuid.UUID

	// User identifier
	UserID *uuid.UUID

	// Media title
	Title string

	// Description media
	Description *string

	// Visibility access level for media (public or private)
	Visibility dtypes.MediaVisibility

	// Number of completed views
	Views uint32

	// Media source creation timestamp
	SourceCreatedAt time.Time
}

func (src *MediaSourceIndex) Copy() *MediaSourceIndex {
	if src == nil {
		return nil
	}

	copy := new(*src)
	copy.UserID = uptr.Copy(src.UserID)
	copy.Description = uptr.Copy(src.Description)

	return copy
}
