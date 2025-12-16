package httpsrv

import (
	"github.com/fasthttp/router"
	apihandlers "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/api"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

// setupAPIV1Routes setup API v1 routes.
func (s *httpServer) setupAPIV1Routes(r *router.Router, handlers *apihandlers.APIHandlers) {
	// Youtube channel
	group := r.Group(httppaths.GroupV1YoutubeChannelClient)
	{
		group.GET(httppaths.PathV1GetYoutubeChannelByID, handlers.V1.GetChannelByID)
	}
}
