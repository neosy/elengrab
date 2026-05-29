package handlers

import (
	"path/filepath"

	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) AssetRobotsHandler(ctx *fasthttp.RequestCtx) {
	fasthttp.ServeFileUncompressed(ctx, filepath.Join(h.assetFolders.Static(), "robots.txt"))
}

func (h *DownloaderHandlers) AssetFabiconHandler(ctx *fasthttp.RequestCtx) {
	fasthttp.ServeFileUncompressed(ctx, filepath.Join(h.assetFolders.Img(), "favicon.ico"))
}
