package handlers

import (
	"strings"

	apierrors "github.com/neosy/elengrab/internal/api/errors"
	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/auth_middleware"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/policy"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) AuthLoginHandler(ctx *fasthttp.RequestCtx) {
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

		ctx.Redirect(httppaths.PathIndex, fasthttp.StatusSeeOther)
		return
	}

	cssPaths, err := uivalues.CssAuthPaths(h.assetFolders.Css())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	jsScripts, err := uivalues.JsAuthPaths(h.assetFolders.Js())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	baseValues := uivalues.NewBaseValues()
	baseValues.Title = uivalues.PageAuthLoginTitle

	dataMap := uivalues.MergeMaps(
		baseValues.ToMap(),
		uivalues.PathValues,
	)
	dataMap[uivalues.CssPathsKey] = cssPaths
	dataMap[uivalues.JsScriptsKey] = jsScripts

	// Load template
	tmpl, err := h.loadPage(uivalues.PageAuthLogin.FileName())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType("text/html; charset=utf-8")

	// Execute template
	if err := tmpl.ExecuteTemplate(ctx, uivalues.PageAuthLogin.Key(), dataMap); err != nil {
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
		ctx.Response.Header.Set("HX-Redirect", httppaths.PathIndex)
		return
	}

	// Get form data
	formLogin := strings.TrimSpace(string(ctx.FormValue(formFieldLoginKey)))
	formPassword := strings.TrimSpace(string(ctx.FormValue(formFieldPasswordKey)))

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

	ctx.SetStatusCode(fasthttp.StatusSeeOther)
	ctx.Response.Header.Set("Cache-Control", "no-store")
	ctx.Response.Header.Set("HX-Redirect", httppaths.PathIndex)
}
