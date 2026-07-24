package watchchunk

import (
	"context"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *MediaWatchChunk) AddChunkQty(ctx context.Context, chunk *ddownload.MediaWatchChunk) error {
	if chunk == nil {
		uc.logger.Warn("Nil pointer in function")
		return apperrors.ErrFuncParamNullPointer
	}

	if err := chunk.Validate(); err != nil {
		return err
	}

	err := uc.chunkRep.AddChunkQty(ctx, chunk)
	if err != nil {
		uc.logger.Warn(
			"Failed to insert record into repository",
			"error", err,
		)
		return errorx.Errorf("failed to insert record: %w", err, exceptionx.ERROR)
	}

	return nil
}

func (uc *MediaWatchChunk) AddChunkQtyBatch(ctx context.Context, chunks []*ddownload.MediaWatchChunk) error {
	if chunks == nil {
		uc.logger.Warn("Nil pointer in function")
		return apperrors.ErrFuncParamNullPointer
	}

	for _, chunk := range chunks {
		if err := chunk.Validate(); err != nil {
			return err
		}
	}

	err := uc.chunkRep.AddChunkQtyBatch(ctx, chunks)
	if err != nil {
		uc.logger.Warn(
			"Failed to insert records into repository",
			"error", err,
		)
		return errorx.Errorf("failed to insert records: %w", err, exceptionx.ERROR)
	}

	return nil
}
