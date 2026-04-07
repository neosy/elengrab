package handlers

import (
	"path/filepath"

	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) RooFilesHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case httppaths.PathRootRobotsTxt:
		fasthttp.ServeFileUncompressed(ctx, filepath.Join(h.assetFolders.Static(), "robots.txt"))
		return
	case httppaths.PathRootFaviconICO:
		fasthttp.ServeFileUncompressed(ctx, filepath.Join(h.assetFolders.Img(), "favicon.ico"))
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.Response.Header.SetContentType("text/plain; charset=utf-8")
	ctx.WriteString("404 Not Found\n")
}
