package handlers

import (
	"fmt"
	"path/filepath"

	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/auth_middleware"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) AuthRegisterHandler(ctx *fasthttp.RequestCtx) {
	// Check if user is already logged in and not a guest
	ctxUser := authmw.UserFromContext(ctx)
	if ctxUser != nil && ctxUser.UserType() != dtypes.UserTypeGuest {
		ctx.Redirect(httppaths.PathIndex, fasthttp.StatusFound)
		return
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType("text/html; charset=utf-8")

	cssPaths, err := uivalues.CssAuthPaths(filepath.Join(h.assetsDir, dirStaticName, dirCssName))
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	jsScripts, err := uivalues.JsAuthScripts(filepath.Join(h.assetsDir, dirStaticName, dirJsName))
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	dataMap := uivalues.MergeMaps(uivalues.IndexValues, uivalues.PathValues)
	dataMap[uivalues.CssPathsKey] = cssPaths
	dataMap[uivalues.JsScriptsKey] = jsScripts

	// Execute template with PageTitle
	if err := h.templates.ExecuteTemplate(ctx, uivalues.AuthRegisterHtmlFileName, dataMap); err != nil {
		nfasthttp.WriteError(ctx, fmt.Errorf("template execution error: %v", err), fasthttp.StatusInternalServerError)
		return
	}
}

func (h *DownloaderHandlers) AuthRegisterSubmitHandler(ctx *fasthttp.RequestCtx) {
	// Check if user is already logged in and not a guest
	ctxUser := authmw.UserFromContext(ctx)
	if ctxUser != nil && ctxUser.UserType() != dtypes.UserTypeGuest {
		ctx.SetStatusCode(fasthttp.StatusSeeOther)
		ctx.Response.Header.Set("HX-Redirect", httppaths.PathIndex)
		return
	}

	// Get form data
	formLogin := string(ctx.FormValue(formFieldLoginKey))
	formPassword := string(ctx.FormValue(formFieldPasswordKey))
	formConfirmPassword := string(ctx.FormValue(formFieldConfirmPasswordKey))

	// Check if passwords match
	if formPassword != formConfirmPassword {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTP("Passwords do not match. Please try again.", fasthttp.StatusBadRequest))
		return
	}

	req := &udto.RegisterUserRequest{
		Login:    formLogin,
		Password: formPassword,
	}

	// Register user and get response token
	resp, err := h.authWeb.RegisterUser(ctx, req)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	if resp == nil || resp.Token == nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTP("a system error. The query returned an empty pointer", fasthttp.StatusInternalServerError))
		return
	}

	// Set session token and expiration time
	authmw.CookieSessionTokenKey.SetValue(ctx, resp.Token.Token, resp.Token.ExpiresAt)

	// Migrate guest data if user was previously a guest
	if ctxUser != nil && ctxUser.UserType() == dtypes.UserTypeGuest {
		err := h.downloader.MigrateGuestData(ctx, ctxUser.UserID, resp.UserID)
		if err == nil {
			h.authWeb.SoftDeleteUser(ctx, ctxUser.UserID)
		}
	}

	ctx.SetStatusCode(fasthttp.StatusSeeOther)
	ctx.Response.Header.Set("HX-Redirect", httppaths.PathIndex)
}
