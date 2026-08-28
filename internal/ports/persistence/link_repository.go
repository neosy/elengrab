package persistence

import (
	"context"

	"github.com/google/uuid"
	dlink "github.com/neosy/elengrab/internal/domain/link"
)

type LinkRepositoryFactory func() LinkRepository

type LinkRepository interface {
	// Insert inserting a record
	Insert(ctx context.Context, link *dlink.Link) error
	// Update updating a record
	Update(ctx context.Context, link *dlink.Link) error
	// Delete soft deleting a record
	SoftDelete(ctx context.Context, linkId uuid.UUID) error
	// Delete hard deleting a record
	HardDelete(ctx context.Context, linkId uuid.UUID) error

	// Find link by linkId
	Find(ctx context.Context, linkId uuid.UUID) (*dlink.Link, error)
	// Exists link by linkId
	Exists(ctx context.Context, linkId uuid.UUID) (bool, error)
	// FindLastByShortCode link by shortCode
	FindLastByShortCode(ctx context.Context, shortCode string) (*dlink.Link, error)
	// ExistsActiveShortCode exists link by active shortCode
	ExistsActiveShortCode(ctx context.Context, shortCode string) (bool, error)

	WithoutDeleted() LinkRepository
}
