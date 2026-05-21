package migrations

import (
	"context"
	"log/slog"

	dlmigration "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_data_migration"
	fileuc "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	pservices "github.com/neosy/elengrab/internal/ports/services"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type usecases struct {
	download          *fileuc.MediaDownload
	downloadMigration *dlmigration.DownloadMigration
	thumbnail         *thumbnail.Thumbnail
}

type services struct {
	downloader pservices.Downloader
	ffmpeg     pservices.FFMpeg
}

type migrations struct {
	logger *slog.Logger

	requiredMigrationList migrationList
	deferredMigrationList migrationList

	dlStorage pstorage.DownloadsStorage

	usecases usecases
	services services
}

func NewMigrations(
	logger *slog.Logger,
	dlStorage pstorage.DownloadsStorage,
	// usecases
	download *fileuc.MediaDownload,
	downloadMigration *dlmigration.DownloadMigration,
	thumbnail *thumbnail.Thumbnail,
	// services
	downloader pservices.Downloader,
	ffmpeg pservices.FFMpeg,
) *migrations {
	migrations := &migrations{
		logger: logger,

		requiredMigrationList: NewMigrationList(),
		deferredMigrationList: NewMigrationList(),

		dlStorage: dlStorage,

		usecases: usecases{
			download:          download,
			downloadMigration: downloadMigration,
			thumbnail:         thumbnail,
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

	err := m.usecases.downloadMigration.Insert(ctx, migration)
	if err != nil {
		return err
	}

	return nil
}
