package apihandlers

import (
	handlers "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/api/v1"
	"github.com/neosy/elengrab/internal/app/usecases"
)

type APIHandlers struct {
	V1 *handlers.V1Handlers
}

func NewAPIHandlers(usecases *usecases.Usecases) *APIHandlers {
	return &APIHandlers{
		V1: handlers.NewV1Handlers(usecases),
	}
}
