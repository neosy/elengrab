package linkclick

import (
	"context"

	"github.com/google/uuid"
)

func (u *LinkClick) CountByLinkId(ctx context.Context, linkID uuid.UUID) (uint16, error) {
	return u.linkClickRepo().CountByLinkId(ctx, linkID)
}
