package handlers

import (
	"errors"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
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
		hash, _ := ctx.UserValue(hashKey).(string)

		// remove query string if there is
		if idx := strings.IndexByte(uri, '?'); idx >= 0 {
			uri = uri[:idx]
		}

		// protection from path traversal
		cleanPath := path.Clean("/" + uri)
		filePath := filepath.Join(subdirPath, filepath.FromSlash(cleanPath))

		staticFile, err := h.readAssetFile(filePath, hash)
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

func (h *StaticHandlers) readAssetFile(filePath string, hash string) (*dtypes.AssetFile, error) {
	if hash != "" {
		assetFile, status, _ := h.assetFileCacheRep.Find(hash)
		if status == memsimple.CacheStatusNegativeHit {
			return nil, errorx.New("file not found", exceptionx.NOT_FOUND)
		}
		if assetFile != nil {
			return assetFile, nil
		}
	}

	data, err := assets.ReadAssetFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			h.assetFileCacheRep.SaveNegative(hash, filePath)
		}
		return nil, err
	}

	file := dtypes.NewAssetFile(filePath, hash, data)

	h.assetFileCacheRep.Save(file)

	return file, nil
}
