package link

import (
	"context"

	"github.com/google/uuid"
)

func (u *Link) SoftDelete(ctx context.Context, linkID uuid.UUID) error {
	return u.linkRepo().SoftDelete(ctx, linkID)
}

func (u *Link) HardDelete(ctx context.Context, linkID uuid.UUID) error {
	return u.linkRepo().HardDelete(ctx, linkID)
}
