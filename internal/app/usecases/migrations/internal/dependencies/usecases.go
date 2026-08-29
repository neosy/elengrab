package dependencies

import (
	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	dlmigration "github.com/neosy/elengrab/internal/app/usecases/migrations/internal/download_data_migration"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
)

type Usecases struct {
	Downloader    downloader.InternalDownloader
	MediaDownload downloader.MediaDownload
	MediaWatch    downloader.MediaWatch

	DownloadMigration *dlmigration.DownloadMigration
	Thumbnail         *thumbnail.Thumbnail
}
