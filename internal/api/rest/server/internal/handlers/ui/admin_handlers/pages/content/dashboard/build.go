package dashboard

import (
	"context"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/pages"
)

func (c *Content) Build(ctx context.Context) (any, error) {
	return pages.AdminContentDashboardData{
		AdminContentData: pages.AdminContentData{
			Title: c.title,
		},
	}, nil
}
