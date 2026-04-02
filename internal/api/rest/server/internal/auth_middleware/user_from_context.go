package authmw

import (
	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/valyala/fasthttp"
)

func UserFromContext(ctx *fasthttp.RequestCtx) *dauth.UserContext {
	ctxUserIface := ctx.UserValue(userKey)
	if ctxUserIface == nil {
		return nil
	}

	userCtx, ok := ctxUserIface.(dauth.UserContext)
	if !ok {
		return nil
	}

	if userCtx.UserID == uuid.Nil {
		return nil
	}

	return &userCtx
}

func EnsureUserFromContext(ctx *fasthttp.RequestCtx) (*dauth.UserContext, error) {
	userCtx := UserFromContext(ctx)
	if userCtx == nil {
		return nil, errorx.NewHTTP("unauthorized", fasthttp.StatusUnauthorized)
	}

	return userCtx, nil
}
