package link

import (
	"context"

	"github.com/google/uuid"
)

func (u *Link) SoftDelete(ctx context.Context, linkID uuid.UUID) error {
	return u.linkRep.SoftDelete(ctx, linkID)
}

func (u *Link) HardDelete(ctx context.Context, linkID uuid.UUID) error {
	return u.linkRep.HardDelete(ctx, linkID)
}
