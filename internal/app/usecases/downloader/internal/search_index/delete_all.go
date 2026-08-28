package searchindex

import (
	"context"

	"github.com/google/uuid"
)

func (uc *SearchIndex) SoftDeleteMediaDownload(ctx context.Context, downloadID uuid.UUID) error {
	return uc.searchIndex.Tx(ctx, func(ctx context.Context) error {
		return uc.searchIndex.SoftDelete(ctx, downloadID)
	})
}

func (uc *SearchIndex) HardDeleteMediaDownload(ctx context.Context, downloadID uuid.UUID) error {
	return uc.searchIndex.Tx(ctx, func(ctx context.Context) error {
		return uc.searchIndex.HardDelete(ctx, downloadID)
	})
}
