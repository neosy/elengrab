package handlers

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/api/v1/mappers"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/api/v1/validators"
	"github.com/neosy/elengrab/internal/app/usecases"
)

type V1Handlers struct {
	mappers    *mappers.Mappers
	validators *validators.Validators
	usecases   *usecases.Usecases
}

func NewV1Handlers(usecases *usecases.Usecases) *V1Handlers {
	return &V1Handlers{
		mappers:    mappers.NewMappers(),
		validators: validators.NewValidators(),
		usecases:   usecases,
	}
}
