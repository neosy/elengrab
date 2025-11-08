package grabh

import "github.com/neosy/elengrab/internal/app/usecases"

type GrabHandlers struct {
	usecases *usecases.Usecases

	// Options
	assetsDir string
}

func NewGrabHandlers(usecases *usecases.Usecases, assetsDir string) *GrabHandlers {
	return &GrabHandlers{
		usecases:  usecases,
		assetsDir: assetsDir,
	}
}
