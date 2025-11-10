package statich

import (
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/paths"
	"github.com/valyala/fasthttp"
)

func (h *StaticHandlers) StaticHandler(ctx *fasthttp.RequestCtx) {
	// Remove /static prefix
	ctx.Request.SetRequestURIBytes(ctx.Path()[len(httppaths.GroupStatic):])

	fs := &fasthttp.FS{
		Root:               h.assetsDir + "/static",
		GenerateIndexPages: false,
	}

	fs.NewRequestHandler()(ctx)
}
