package downloader

import (
	"mime"
	"strings"

	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/clientcap"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/middleware/auth"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) AuthLoginPageHandler(ctx *fasthttp.RequestCtx) {
	if err := nfasthttp.EnforceHTTPS(ctx); err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	// Check if user is already logged in and not a guest
	ctxUser := policy.ResolveUser(ctx)
	if ctxUser != nil && ctxUser.UserType() != dtypes.UserTypeGuest {
		ctx.Response.Header.Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		ctx.Response.Header.Set("Pragma", "no-cache")
		ctx.Response.Header.Set("Expires", "0")

		ctx.Redirect(httppaths.IndexPath, fasthttp.StatusSeeOther)
		return
	}

	cssPaths, err := h.assetPaths.AuthPageCssPaths()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	caps := clientcap.Detect(string(ctx.UserAgent()))

	jsScripts, err := h.assetPaths.AuthPageJsPaths(caps.IsLegacyWebKit)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	baseValues := pages.NewBaseValues()
	baseValues.Title = pages.PageAuthLoginTitle

	redirectPath := string(ctx.QueryArgs().Peek(redirectKey))
	if redirectPath == "" {
		redirectPath = httppaths.IndexPath
	}

	pageData := pages.AuthLoginPageData{
		BaseValues: baseValues,
		BasePaths:  paths.NewHttpPaths(),
		Paths: pages.PagePaths{
			Css:       cssPaths,
			JsScripts: jsScripts,
		},
		Values: pages.AuthLoginPageValues{
			Redirect: redirectPath,
		},
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType(mime.TypeByExtension(".html"))

	// Execute template
	if err := h.templates.Pages[pages.AuthLoginPage.Key()].ExecuteTemplate(ctx, pages.AuthLoginPage.Key(), pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}

func (h *DownloaderHandlers) AuthLoginSubmitHandler(ctx *fasthttp.RequestCtx) {
	if err := nfasthttp.EnforceHTTPS(ctx); err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	// Check if user is already logged in and not a guest
	ctxUser := policy.ResolveUser(ctx)
	if ctxUser != nil && ctxUser.UserType() != dtypes.UserTypeGuest {
		ctx.SetStatusCode(fasthttp.StatusSeeOther)
		ctx.Response.Header.Set("HX-Redirect", httppaths.IndexPath)
		return
	}

	// Get form data
	formLogin := strings.TrimSpace(string(ctx.FormValue(formFieldLoginKey)))
	formPassword := strings.TrimSpace(string(ctx.FormValue(formFieldPasswordKey)))
	formRedirect := strings.TrimSpace(string(ctx.FormValue(formFieldRedirectKey)))

	req := &udto.AuthUserRequest{
		Login:    formLogin,
		Password: formPassword,
	}

	// Register user and get response token
	resp, err := h.authWeb.LoginUser(ctx, req)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	if resp == nil || resp.Token == nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrQueryReturnedEmptyResult)
		return
	}

	// Set session token and expiration time
	authmw.CookieSessionTokenKey.SetValueWithSecure(ctx, resp.Token.Token, authmw.WithExpiresAt(resp.Token.ExpiresAt))

	// Migrate guest data if user was previously a guest
	if ctxUser != nil && ctxUser.UserType() == dtypes.UserTypeGuest {
		err := h.downloader.MigrateGuestData(ctx, ctxUser.UserID, resp.UserID)
		if err == nil {
			h.authWeb.SoftDeleteUser(ctx, ctxUser.UserID)
		}
	}

	redirectPath := formRedirect
	if redirectPath == "" {
		redirectPath = httppaths.IndexPath
	}

	ctx.SetStatusCode(fasthttp.StatusSeeOther)
	ctx.Response.Header.Set("Cache-Control", "no-store")
	ctx.Response.Header.Set("HX-Redirect", redirectPath)
}
