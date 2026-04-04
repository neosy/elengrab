package persistence

import (
	"context"

	"github.com/google/uuid"
	dlink "github.com/neosy/elengrab/internal/domain/link"
)

type LinkClickRepository interface {
	// Insert inserting a record
	Insert(ctx context.Context, linkClick *dlink.LinkClick) error

	// Find linkClick by linkClickId
	Find(ctx context.Context, linkClickId uuid.UUID) (*dlink.LinkClick, error)
	// CountByLinkId returns the number of records associated with the given linkId.
	CountByLinkId(ctx context.Context, linkId uuid.UUID) (uint16, error)
}
