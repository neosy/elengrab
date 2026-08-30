package link

import (
	"context"

	"github.com/google/uuid"
)

func (u *link) SoftDelete(ctx context.Context, linkID uuid.UUID) error {
	return u.link.SoftDelete(ctx, linkID)
}
