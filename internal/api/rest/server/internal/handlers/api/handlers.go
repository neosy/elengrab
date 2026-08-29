package api

import (
	handlers "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/api/v1"
	"github.com/neosy/elengrab/internal/app/usecases"
)

type Handlers struct {
	V1 *handlers.V1Handlers
}

func NewHandlers(usecases *usecases.Usecases) *Handlers {
	return &Handlers{
		V1: handlers.NewV1Handlers(usecases.DownloaderAPI),
	}
}
