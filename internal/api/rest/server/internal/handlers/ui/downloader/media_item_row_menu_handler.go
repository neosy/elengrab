package downloader

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/menu"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) MediaItemRowMenuHandler(ctx *fasthttp.RequestCtx) {
	authCtx := policy.ResolveUserOrAnonym(ctx)

	downloadID, err := h.extractDownloadID(ctx)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	downloadInfo, err := h.downloader.GetDownloadInfo(ctx, authCtx, downloadID)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	if ctx.IsHead() {
		fasthttpx.WriteHTMLHeadResponse(ctx)
		return
	}

	var errorMessage string
	if downloadInfo.MediaDownload != nil && downloadInfo.MediaDownload.ErrorMessage != nil {
		errorMessage = *downloadInfo.MediaDownload.ErrorMessage
	}

	link, _ := h.linkWeb.ResolveURL(ctx, h.buildMediaWatchURL(ctx, downloadInfo.DownloadID))

	extraData := make(map[string]any)
	extraData[items.MenuActionsKey] = menu.RowMenuActions(
		map[string]string{
			menu.RowMenuActionItemIDKey: idcodec.EncodeUUIDBase64URL(downloadID),
			menu.RowMenuActionURLKey:    downloadInfo.MediaURL,
		},
		downloadInfo.Status,
		downloadInfo.MediaDownload.SourceType,
		downloadInfo.HasEditAccess,
		downloadInfo.HasDeleteAccess,
		menu.WithErrorText(errorMessage),
		menu.WithShareLink(link),
		menu.WithMetadata(downloadInfo.MediaInfo != nil),
		menu.WithVisibility(downloadInfo.Visibility),
	)

	pageData := pages.PageFragmentData{
		BasePaths: paths.NewHttpPaths(),
		Extra:     extraData,
	}

	// Execute template
	err = h.templates.Base.ExecuteTemplate(ctx, components.MenuContentKey, pageData)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, errInternal(err))
	}
}
