package dlmigration

import (
	"context"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (c *DownloadMigration) Insert(ctx context.Context, migration *ddownload.DataMigration) error {
	if migration == nil {
		c.logger.Warn("Nil pointer in function")
		return apperrors.ErrFuncParamNullPointer
	}

	err := c.dataMigrationRepo().Insert(ctx, migration)
	if err != nil {
		c.logger.Warn(
			"Failed to insert record into repository",
			"migrationID", migration.MigrationID,
			"error", err,
		)
		return err
	}

	return nil
}
