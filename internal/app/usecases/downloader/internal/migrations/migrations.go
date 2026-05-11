package migrations

import (
	"context"
	"log/slog"

	dlmigration "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_data_migration"
	fileuc "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/file"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	pservices "github.com/neosy/elengrab/internal/ports/services"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type usecases struct {
	media          *fileuc.File
	mediaMigration *dlmigration.DataMigration
	thumbnail      *thumbnail.Thumbnail
}

type services struct {
	downloader pservices.Downloader
	ffmpeg     pservices.FFMpeg
}

type migrations struct {
	logger *slog.Logger

	migrationIDs migrationIDMap

	dlStorage pstorage.DownloadsStorage

	usecases usecases
	services services
}

func NewMigrations(
	logger *slog.Logger,
	dlStorage pstorage.DownloadsStorage,
	// usecases
	media *fileuc.File,
	mediaMigration *dlmigration.DataMigration,
	thumbnail *thumbnail.Thumbnail,
	// services
	downloader pservices.Downloader,
	ffmpeg pservices.FFMpeg,
) *migrations {
	migrations := &migrations{
		logger: logger,

		migrationIDs: make(migrationIDMap),

		dlStorage: dlStorage,

		usecases: usecases{
			media:          media,
			mediaMigration: mediaMigration,
			thumbnail:      thumbnail,
		},

		services: services{
			downloader: downloader,
			ffmpeg:     ffmpeg,
		},
	}

	migrations.initMigrations()

	return migrations
}

func (m *migrations) markMigration(ctx context.Context, migrationID string) error {
	migration := &ddownload.DataMigration{
		MigrationID: migrationID,
	}

	err := m.usecases.mediaMigration.Insert(ctx, migration)
	if err != nil {
		return err
	}

	return nil
}
