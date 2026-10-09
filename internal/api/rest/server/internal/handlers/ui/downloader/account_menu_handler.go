package downloader

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/menu"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/stringx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) AccountMenuHandler(ctx *fasthttp.RequestCtx) {
	authCtx, err := policy.EnsureUser(ctx)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	extraData := make(map[string]any)
	extraData[items.UserLoginKey] = stringx.Capitalize(authCtx.Login)
	extraData[items.UserEmailKey] = authCtx.Email
	extraData[items.MenuActionsKey] = menu.NewAccountMenuActions(authCtx)

	type pageFragmentData struct {
		pages.PageFragmentData
		pages.UserAvatar
	}

	pageData := pageFragmentData{
		BasePaths: paths.NewHttpPaths(),
		Extra:     extraData,

		UserAvatar: pages.UserAvatar{
			Icon: icons.UserAvatarIconByType(authCtx.UserType()).FileRaw(),
		},
	}

	// Execute template
	if err := h.templates.Base.ExecuteTemplate(ctx, components.AccountMenuContentKey, pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}
