package persistence

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type DownloadDataMigrationRepository interface {
	Insert(ctx context.Context, migration *ddownload.DataMigration) error
	Find(ctx context.Context, migrationID string) (*ddownload.DataMigration, error)
	Exists(ctx context.Context, migrationID string) (bool, error)
	Tx(ctx context.Context, fn func(ctx context.Context) error) error
}
