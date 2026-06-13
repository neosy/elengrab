package static

import (
	"mime"
	"path"
	"path/filepath"
	"strings"

	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *StaticHandlers) newAssetHandler(subdirPath string) func(*fasthttp.RequestCtx) {
	writeResponse := func(ctx *fasthttp.RequestCtx, fileName string, data []byte) {
		contentType := mime.TypeByExtension(filepath.Ext(fileName))
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		ctx.SetContentType(contentType)
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Write(data)
	}

	handler := func(ctx *fasthttp.RequestCtx) {
		uri := string(ctx.RequestURI())

		// remove query string if there is
		if idx := strings.IndexByte(uri, '?'); idx >= 0 {
			uri = uri[:idx]
		}

		// protection from path traversal
		cleanPath := path.Clean("/" + uri)
		filePath := filepath.Join(subdirPath, filepath.FromSlash(cleanPath))

		assetFile, err := h.assets.ReadAssetFile(filePath)
		if err != nil {
			nfasthttp.WriteErrorx(ctx, err)
			return
		}

		writeResponse(ctx, assetFile.FileName, assetFile.Raw)
	}

	return handler
}
