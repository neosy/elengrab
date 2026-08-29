package apiv1

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/api/v1/mappers"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/api/v1/validators"
	"github.com/neosy/elengrab/internal/app/usecases/downloader"
)

type V1Handlers struct {
	mappers    *mappers.Mappers
	validators *validators.Validators
	downloader downloader.DownloaderAPI
}

func NewV1Handlers(downloader downloader.DownloaderAPI) *V1Handlers {
	return &V1Handlers{
		mappers:    mappers.NewMappers(),
		validators: validators.NewValidators(),
		downloader: downloader,
	}
}
