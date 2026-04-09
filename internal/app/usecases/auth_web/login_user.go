package authweb

import (
	"context"
	"strings"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/valyala/fasthttp"
)

func (a *AuthWeb) LoginUser(
	ctx context.Context,
	req *dto.AuthUserRequest,
) (*dto.AuthUserResponse, error) {
	if req == nil {
		return nil, errorx.New("function parameter is nil", exceptionx.ERROR)
	}

	if strings.TrimSpace(req.Login) == "" {
		return nil, errorx.NewHTTP("login is required", fasthttp.StatusBadRequest)
	}

	if strings.TrimSpace(req.Password) == "" {
		return nil, errorx.NewHTTP("password is required", fasthttp.StatusBadRequest)
	}

	resp, err := a.auth.AuthenticateUser(ctx, req)
	if err != nil {
		exception := errorx.OuterException(err)
		if exception == nil {
			return nil, err
		}
		if exception.Num() == exceptionx.NOT_FOUND.Num() || exception.Num() == exceptionx.UNAUTHORIZED.Num() {
			return nil, errorx.NewHTTP("invalid login or password", fasthttp.StatusUnauthorized)
		}
		return nil, err
	}

	return resp, nil
}
