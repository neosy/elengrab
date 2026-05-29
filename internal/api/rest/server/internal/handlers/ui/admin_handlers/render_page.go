package handlers

import (
	"html/template"
	"mime"
	"strings"

	navmenu "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin_handlers/nav_menu"
	adminpages "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin_handlers/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/paths"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *AdminHandlers) renderPage(ctx *fasthttp.RequestCtx, ctxUser *dauth.UserContext) {
	cssPaths, err := paths.AdminCssPaths(h.assetFolders.Css())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	jsScripts, err := paths.AdminPageJsPaths(h.assetFolders.Js())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	jsImportJSON, err := paths.AdminPageJsImportJSON(h.assetFolders.Js())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	baseValues := pages.NewBaseValues()

	extraData := make(map[string]any)
	extraData[items.UserAvatarIconKey] = icons.UserAvatarIconByType(ctxUser.UserType()).FileRaw()

	page := adminpages.PageByURI(string(ctx.RequestURI()), httppaths.AdminGroup)

	pageData := pages.AdminPageData{
		BasePaths:  paths.NewPaths(),
		BaseValues: baseValues,
		Paths: pages.PagePaths{
			Css:             cssPaths,
			JsScripts:       jsScripts,
			JsImportMapJSON: template.HTML(jsImportJSON),
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

	// Load template
	tmpl, err := h.loadPageTemplate(pages.AdminPage.FileName())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errTemplateInternal(err))
		return
	}

	// Execute template with PageTitle
	if err := tmpl.ExecuteTemplate(ctx, pages.AdminPage.Key(), pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errTemplateInternal(err))
		return
	}
}
