package statich

import (
	handlers "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/static/handlers"
	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
)

type StaticHandlers struct {
	Static *handlers.StaticHandlers
}

func NewStaticHandlers(
	assetsDir string,

	// usecases
	thumbnail *thumbnail.Thumbnail,
	downloader *downloader.Downloader,
) *StaticHandlers {
	return &StaticHandlers{
		Static: handlers.NewStaticHandlers(assetsDir, thumbnail, downloader),
	}
}
