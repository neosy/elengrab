package dashboard

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin/mappers"
	adminuc "github.com/neosy/elengrab/internal/app/usecases/admin"
)

type Content struct {
	mappers *mappers.Mappers

	admin *adminuc.Admin

	title string
}

func NewContent(mappers *mappers.Mappers, admin *adminuc.Admin, title string) *Content {
	return &Content{
		mappers: mappers,
		admin:   admin,
		title:   title,
	}
}
