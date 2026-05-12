package dlmigration

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (c *DownloadMigration) Find(ctx context.Context, migrationID string) (*ddownload.DataMigration, error) {
	migration, err := c.dataMigrationRep.Find(ctx, migrationID)
	if err != nil {
		c.logger.Warn("Failed to find migration", "error", err)
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}
	return migration, nil
}

func (c *DownloadMigration) GetByMigrationID(ctx context.Context, migrationID string) (*ddownload.DataMigration, error) {
	migration, err := c.Find(ctx, migrationID)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if migration == nil {
		c.logger.Warn("Thumbnail not found", "migrationID", migrationID)
		return nil, ierrors.ErrThumbnailNotFound
	}

	return migration, nil
}

func (c *DownloadMigration) Exists(ctx context.Context, migrationID string) (bool, error) {
	exists, err := c.dataMigrationRep.Exists(ctx, migrationID)
	if err != nil {
		c.logger.Warn("Failed to exists migration", "error", err)
		return false, errorx.NewFromError(err, exceptionx.ERROR)
	}

	return exists, nil
}
