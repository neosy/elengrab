package handlers

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/menu"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) SettingsMenuHandler(ctx *fasthttp.RequestCtx) {
	extraData := make(map[string]any)
	extraData[items.MenuActionsKey] = menu.SettingsMenuActions()

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
	if err := tmpl.ExecuteTemplate(ctx, components.MenuContentKey, pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}
