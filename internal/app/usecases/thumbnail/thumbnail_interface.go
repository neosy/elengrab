package thumbnail

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

type Thumbnail interface {
	Create(ctx context.Context, req *dto.CreateThumbnailRequest) (uuid.UUID, error)
	Delete(ctx context.Context, thumbID uuid.UUID) error
	DeleteBatch(ctx context.Context, thumbIDs []uuid.UUID) error

	// FindByThumbID retrieves thumbnail information by its ID without loading the raw image data.
	// Returns nil if the thumbnail does not exist.
	FindByThumbID(ctx context.Context, thumbID uuid.UUID) (*dmedia.Thumbnail, error)

	// GetByThumbID retrieves thumbnail information by its ID without loading the raw image data.
	GetByThumbID(ctx context.Context, thumbID uuid.UUID) (*dmedia.Thumbnail, error)

	// LoadByThumbID retrieves thumbnail metadata by ID and loads its file data from storage.
	LoadByThumbID(ctx context.Context, thumbID uuid.UUID) (*dmedia.Thumbnail, error)
}
