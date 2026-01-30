package handlers

import (
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/valyala/fasthttp"
)

type StaticHandlers struct {
	assetsDir string

	cssHandler  fasthttp.RequestHandler
	imgHandler  fasthttp.RequestHandler
	iconHandler fasthttp.RequestHandler
	jsHandler   fasthttp.RequestHandler
	pwaHandler  fasthttp.RequestHandler
}

func NewStaticHandlers(assetsDir string) *StaticHandlers {

	h := &StaticHandlers{
		assetsDir: assetsDir,
	}

	h.cssHandler = h.newFSHandler("css", "css")
	h.imgHandler = h.newFSHandler("img", "img")
	h.iconHandler = h.newFSHandler("img/icons", "icon")
	h.jsHandler = h.newFSHandler("js", "js")
	h.pwaHandler = h.newFSHandler("pwa", "pwa")

	return h
}

func (h *StaticHandlers) newFSHandler(name string, htmlPathName string) fasthttp.RequestHandler {
	fs := &fasthttp.FS{
		Root:               h.assetsDir + "/static/" + name,
		GenerateIndexPages: false,
	}

	handler := fs.NewRequestHandler()

	return func(ctx *fasthttp.RequestCtx) {
		path := httppaths.GroupStatic + "/" + htmlPathName
		ctx.Request.SetRequestURIBytes(
			ctx.Path()[len(path):],
		)
		handler(ctx)
	}
}
