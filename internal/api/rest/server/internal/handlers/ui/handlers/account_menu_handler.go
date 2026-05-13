package handlers

import (
	"html/template"
	"path/filepath"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/policy"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) AccountMenuHandler(ctx *fasthttp.RequestCtx) {
	ctxUser, err := policy.EnsureUser(ctx)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")

	extraData := make(map[string]any)
	extraData[uivalues.UserAvatarIconKey] = template.HTML(
		uivalues.IconFileRawByKey(uivalues.UserAvatarKeyByType(ctxUser.UserType()), iconsDir))
	extraData[uivalues.UserLoginKey] = capitalize(ctxUser.Login)
	extraData[uivalues.UserEmailKey] = ctxUser.Email
	extraData[uivalues.AccountMenuActionsKey] = uivalues.AccountMenuActions(iconsDir)

	pageData := uivalues.PageFragmentData{
		BasePaths: uivalues.NewBasePaths(),
		Extra:     extraData,
	}

	// Load template
	tmpl, err := h.templates.Clone()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	// Execute template
	if err := tmpl.ExecuteTemplate(ctx, uivalues.ComponentAccountMenuContentKey, pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}
