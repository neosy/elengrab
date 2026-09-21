package searchindex

import (
	"context"

	"github.com/google/uuid"
)

func (uc *SearchIndex) SoftDeleteMediaDownload(ctx context.Context, downloadID uuid.UUID) error {
	return uc.sourceIndex.Tx(ctx, func(ctx context.Context) error {
		return uc.sourceIndex.SoftDelete(ctx, downloadID)
	})
}

func (uc *SearchIndex) HardDeleteMediaDownload(ctx context.Context, downloadID uuid.UUID) error {
	return uc.sourceIndex.Tx(ctx, func(ctx context.Context) error {
		return uc.sourceIndex.HardDelete(ctx, downloadID)
	})
}
