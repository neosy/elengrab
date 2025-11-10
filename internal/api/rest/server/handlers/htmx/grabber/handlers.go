package grabberh

import "github.com/neosy/elengrab/internal/app/usecases"

type GrabberHandlers struct {
	usecases *usecases.Usecases

	// Options
	assetsDir string
}

func NewGrabberHandlers(usecases *usecases.Usecases, assetsDir string) *GrabberHandlers {
	return &GrabberHandlers{
		usecases:  usecases,
		assetsDir: assetsDir,
	}
}
