package routes

import (
	apiv1 "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/api/v1"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

// registerAPI register API routes.
func (r *routes) registerAPI(v1 *apiv1.V1Handlers) {
	// Youtube channel
	group := r.router.Group(httppaths.APIV1YoutubeChannelClientGroup)
	{
		group.GET(httppaths.APIV1GetYoutubeChannelByIDPath, v1.GetChannelByID)
	}
}
