package link

import (
	"context"

	"uuid"
)

func (u *link) SoftDelete(ctx context.Context, linkID uuid.UUID) error {
	return u.link.SoftDelete(ctx, linkID)
}
