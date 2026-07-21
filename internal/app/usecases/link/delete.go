package link

import (
	"context"

	"github.com/google/uuid"
)

func (u *Link) SoftDelete(ctx context.Context, linkID uuid.UUID) error {
	return u.link.SoftDelete(ctx, linkID)
}
