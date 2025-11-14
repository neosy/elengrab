package grabberh

import (
	"github.com/neosy/elengrab/internal/api/rest/server/handlers/htmx/grabber/validators"
	"github.com/neosy/elengrab/internal/app/usecases"
)

type GrabberHandlers struct {
	validators *validators.Validators
	usecases   *usecases.Usecases

	// Options
	assetsDir string
}

func NewGrabberHandlers(usecases *usecases.Usecases, assetsDir string) *GrabberHandlers {
	return &GrabberHandlers{
		validators: validators.NewValidators(),
		usecases:   usecases,
		assetsDir:  assetsDir,
	}
}
