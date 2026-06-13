package static

import (
	"errors"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
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

		staticFile, err := h.assets.ReadAssetFile(filePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				nfasthttp.WriteErrorx(ctx, errorx.New("file not found", exceptionx.NOT_FOUND))
				return
			}
			nfasthttp.WriteErrorx(ctx, err)
			return
		}

		writeResponse(ctx, staticFile.FileName, staticFile.Raw)
	}

	return handler
}
