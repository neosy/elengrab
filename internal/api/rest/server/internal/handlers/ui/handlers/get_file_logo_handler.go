package handlers

import (
	"path/filepath"

	"github.com/google/uuid"
	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/auth_middleware"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) GetFileLogoHandler(ctx *fasthttp.RequestCtx) {
	ctxUser := authmw.UserFromContext(ctx)
	if ctxUser == nil {
		// anonymous
		ctxUser = dauth.UserContextAnonymous()
	}

	fileIdStr := ctx.UserValue(fileIdKey).(string)
	if fileIdStr == "" {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("fileId is required", fasthttp.StatusBadRequest))
		return
	}

	fileID, err := uuid.Parse(fileIdStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("fileId is incorrect", fasthttp.StatusBadRequest, err))
		return
	}

	iconImage, err := h.downloader.GetIcon(ctx, *ctxUser, fileID)

	if err == nil && iconImage != nil && len(iconImage.Raw) > 0 {
		ctx.SetContentType(httpx.ContentTypeByExt(iconImage.Format))
		ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
		ctx.SetBody(iconImage.Raw)
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")
	defaultAvatarSVG := uivalues.IconFileRawByKey(uivalues.MediaDefaultIconNameKey, iconsDir)

	ctx.SetContentType("image/svg+xml")
	ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
	ctx.SetBody([]byte(defaultAvatarSVG))
	ctx.SetStatusCode(fasthttp.StatusOK)
}
