package content

import (
	"context"
	"fmt"

	adminuc "github.com/neosy/elengrab/internal/app/usecases/admin"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin/mappers"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin/pages/content/dashboard"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin/pages/content/users"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin/pages/page"
)

type Content interface {
	Build(ctx context.Context) (any, error)
}

func NewContent(mappers *mappers.Mappers, admin *adminuc.Admin, p page.Page) (Content, error) {
	switch p.Name {
	case page.PageNameDashboard:
		return dashboard.NewContent(mappers, admin, p.Title), nil
	case page.PageNameUsers:
		return users.NewContent(mappers, admin, p.Title), nil
	}

	return nil, fmt.Errorf("unsupported page name: %q", p.Name)
}
