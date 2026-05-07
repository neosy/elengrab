package handlers

import (
	"path/filepath"

	"github.com/google/uuid"
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/policy"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) RowMenuHandler(ctx *fasthttp.RequestCtx) {
	ctxUser := policy.ResolveUserOrAnonym(ctx)

	fileIdStr, ok := ctx.UserValue(fileIdKey).(string)
	if !ok || fileIdStr == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrFileIdIsRequired)
		return
	}

	fileId, err := uuid.Parse(fileIdStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrFileIdIsIncorrect.Wrap(err))
		return
	}

	fileResp, err := h.downloader.GetFileInfo(ctx, *ctxUser, fileId)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")

	dataMap := uivalues.MergeMaps(
		uivalues.PathValues,
	)
	dataMap[uivalues.RowMenuActionsKey] = uivalues.RowMenuActions(
		iconsDir,
		map[string]string{
			uivalues.RowMenuActionFileIdKey: fileId.String(),
			uivalues.RowMenuActionURLKey:    fileResp.MediaUrl,
		},
		fileResp.Status == dtypes.FileStatusDone,
	)

	// Load template
	tmpl, err := h.templates.Clone()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	// Execute template
	if err := tmpl.ExecuteTemplate(ctx, uivalues.ComponentRowMenuContentKey, dataMap); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}
