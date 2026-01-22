package statich

import (
	handlers "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/static/handlers"
)

type StaticHandlers struct {
	Static *handlers.StaticHandlers
}

func NewStaticHandlers(
	assetsDir string,
) *StaticHandlers {
	return &StaticHandlers{
		Static: handlers.NewStaticHandlers(assetsDir),
	}
}
