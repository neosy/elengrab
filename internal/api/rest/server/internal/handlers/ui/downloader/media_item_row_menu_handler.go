package downloader

import (
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/menu"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) MediaItemRowMenuHandler(ctx *fasthttp.RequestCtx) {
	authCtx := policy.ResolveUserOrAnonym(ctx)

	downloadIDStr, ok := ctx.UserValue(downloadIDKey).(string)
	if !ok || downloadIDStr == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsRequired)
		return
	}

	downloadID, err := idcodec.DecodeUUIDBase64URL(downloadIDStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsIncorrect.Wrap(err))
		return
	}

	downloadResp, err := h.downloader.GetDownloadInfo(ctx, authCtx, downloadID)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	extraData := make(map[string]any)
	extraData[items.MenuActionsKey] = menu.RowMenuActions(
		map[string]string{
			menu.RowMenuActionItemIDKey: idcodec.EncodeUUIDBase64URL(downloadID),
			menu.RowMenuActionURLKey:    downloadResp.MediaURL,
		},
		downloadResp.Status,
		downloadResp.HasWriteAccess,
	)

	pageData := pages.PageFragmentData{
		BasePaths: paths.NewHttpPaths(),
		Extra:     extraData,
	}

	// Load template
	tmpl, err := h.templates.Clone()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	// Execute template
	if err := tmpl.ExecuteTemplate(ctx, components.MenuContentKey, pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}
