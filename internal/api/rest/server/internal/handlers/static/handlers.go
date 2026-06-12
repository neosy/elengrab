package statich

import (
	handlers "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/static/handlers"
	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type StaticHandlers struct {
	Static *handlers.StaticHandlers
}

func NewStaticHandlers(
	assetsDir string,

	// caches
	assetFileCacheRepository persistence.AssetFileCacheRepository,

	// usecases
	thumbnail *thumbnail.Thumbnail,
	downloader *downloader.Downloader,
) *StaticHandlers {
	return &StaticHandlers{
		Static: handlers.NewStaticHandlers(assetsDir, assetFileCacheRepository, thumbnail, downloader),
	}
}
