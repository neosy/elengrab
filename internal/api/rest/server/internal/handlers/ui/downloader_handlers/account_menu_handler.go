package handlers

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/menu"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/paths"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/policy"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/stringx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) AccountMenuHandler(ctx *fasthttp.RequestCtx) {
	ctxUser, err := policy.EnsureUser(ctx)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	extraData := make(map[string]any)
	extraData[items.UserAvatarIconKey] = icons.UserAvatarIconByType(ctxUser.UserType()).FileRaw()
	extraData[items.UserLoginKey] = stringx.Capitalize(ctxUser.Login)
	extraData[items.UserEmailKey] = ctxUser.Email
	extraData[items.AccountMenuActionsKey] = menu.NewAccountMenuActions(ctxUser)

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
	if err := tmpl.ExecuteTemplate(ctx, components.AccountMenuContentKey, pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}
