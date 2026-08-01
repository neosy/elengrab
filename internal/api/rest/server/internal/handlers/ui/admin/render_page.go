package admin

import (
	"html/template"
	"mime"
	"strings"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/clientcap"
	navmenu "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin/nav_menu"
	adminpages "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *AdminHandlers) renderPage(ctx *fasthttp.RequestCtx, ctxUser dauth.AuthContext) {
	cssPaths, err := h.assetPaths.AdminPageCssPaths()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	caps := clientcap.Detect(string(ctx.UserAgent()))

	jsScripts, err := h.assetPaths.AdminPageJsPaths(caps.IsLegacyWebKit)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	baseValues := pages.NewBaseValues()

	extraData := make(map[string]any)
	extraData[items.UserAvatarIconKey] = icons.UserAvatarIconByType(ctxUser.UserType()).FileRaw()

	page := adminpages.PageByURI(string(ctx.RequestURI()), httppaths.AdminGroup)

	pageData := pages.AdminPageData{
		BasePaths:  paths.NewHttpPaths(),
		BaseValues: baseValues,
		Paths: pages.PagePaths{
			Css:       cssPaths,
			JsScripts: jsScripts,
		},
		Values: pages.AdminPageValues{
			PageName:            page.Name.String(),
			PageTitle:           page.Title,
			IsPageLogoSymbol:    !page.HasIcon(),
			PageLogoHTML:        page.LogoHTML(),
			NavMenu:             navmenu.NewMenuItems(page.Name),
			ContentTemplateName: page.ContentTemplateName,
		},
		Extra: extraData,
	}

	buf, err := h.renderContent(ctx, page)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	pageData.Values.ContentHTML = template.HTML(strings.TrimSpace(buf.String()))

	// Set content type so browser renders HTML properly
	ctx.SetContentType(mime.TypeByExtension(".html"))

	// Execute template with PageTitle
	if err := h.templates.Pages[pages.AdminPage.Key()].ExecuteTemplate(ctx, pages.AdminPage.Key(), pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errTemplateInternal(err))
		return
	}
}
