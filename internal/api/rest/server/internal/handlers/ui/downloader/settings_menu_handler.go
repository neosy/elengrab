package downloader

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/menu"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) SettingsMenuHandler(ctx *fasthttp.RequestCtx) {
	extraData := make(map[string]any)
	extraData[items.MenuActionsKey] = menu.SettingsMenuActions()

	pageData := pages.PageFragmentData{
		BasePaths: paths.NewHttpPaths(),
		Extra:     extraData,
	}

	// Execute template
	if err := h.templates.Base.ExecuteTemplate(ctx, components.MenuContentKey, pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}
