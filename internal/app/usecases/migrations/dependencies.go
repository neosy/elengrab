package migrations

import (
	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	"github.com/neosy/elengrab/internal/app/usecases/migrations/internal/dependencies"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type Usecases struct {
	Downloader downloader.InternalDownloader
	Thumbnail  thumbnail.Thumbnail
}
type Services = dependencies.Services

type Dependencies struct {
	DownloadsStorage pstorage.DownloadsStorage

	Usecases Usecases
	Services Services
}
