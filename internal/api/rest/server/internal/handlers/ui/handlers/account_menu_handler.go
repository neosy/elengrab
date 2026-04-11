package handlers

import (
	"html/template"
	"path/filepath"

	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/auth_middleware"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) AccountMenuHandler(ctx *fasthttp.RequestCtx) {
	ctxUser, err := authmw.EnsureUserFromContext(ctx)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")

	dataMap := uivalues.MergeMaps(
		uivalues.PathValues,
	)
	dataMap[uivalues.UserAvatarIconKey] = template.HTML(
		uivalues.IconFileRawByKey(uivalues.UserAvatarKeyByType(ctxUser.UserType()), iconsDir))
	dataMap[uivalues.UserLoginKey] = capitalize(ctxUser.Login)
	dataMap[uivalues.UserEmailKey] = ctxUser.Email
	dataMap[uivalues.AccountMenuActionsKey] = uivalues.AccountMenuActions(iconsDir)

	// Load template
	tmpl, err := h.templates.Clone()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	// Execute template
	if err := tmpl.ExecuteTemplate(ctx, uivalues.ComponentAccountMenuContentKey, dataMap); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}
