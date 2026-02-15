package handlers

import (
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) GetFileLogoHandler(ctx *fasthttp.RequestCtx) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetBodyString(fmt.Sprintf("Authorization error: %v", err))
		return
	}

	fileIdStr := ctx.UserValue(fileIdKey).(string)
	if fileIdStr == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("FileId is required")
		return
	}

	fileID, err := uuid.Parse(fileIdStr)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("FileId is incorrect")
		return
	}

	iconImage, err := h.usecases.Downloader.GetIcon(ctx, userID, fileID)

	if err == nil && iconImage != nil && len(iconImage.Raw) > 0 {
		ctx.SetContentType(h.mappers.MapImageExtToContentType(iconImage.Format))
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
