package handlers

import (
	"github.com/google/uuid"
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/menu"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/paths"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/policy"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) MediaItemRowMenuHandler(ctx *fasthttp.RequestCtx) {
	ctxUser := policy.ResolveUserOrAnonym(ctx)

	downloadIDStr, ok := ctx.UserValue(downloadIDKey).(string)
	if !ok || downloadIDStr == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsRequired)
		return
	}

	downloadID, err := uuid.Parse(downloadIDStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsIncorrect.Wrap(err))
		return
	}

	downloadResp, err := h.downloader.GetDownloadInfo(ctx, *ctxUser, downloadID)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	extraData := make(map[string]any)
	extraData[items.RowMenuActionsKey] = menu.RowMenuActions(
		h.assetFolders.Icons(),
		map[string]string{
			menu.RowMenuActionItemIDKey: downloadID.String(),
			menu.RowMenuActionURLKey:    downloadResp.MediaURL,
		},
		downloadResp.Status == dtypes.MediaDownloadStatusDone,
	)

	pageData := pages.PageFragmentData{
		BasePaths: paths.NewPaths(),
		Extra:     extraData,
	}

	// Load template
	tmpl, err := h.templates.Clone()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	// Execute template
	if err := tmpl.ExecuteTemplate(ctx, components.RowMenuContentKey, pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}
