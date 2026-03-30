package authweb

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/valyala/fasthttp"
)

func (a *AuthWeb) RegisterUser(
	ctx context.Context,
	req *dto.RegisterUserRequest,
) (*dto.AuthUserResponse, error) {
	if req == nil {
		return nil, errorx.New("function parameter is nil", exceptionx.ERROR)
	}

	if req.Login == "" {
		return nil, errorx.NewHTTP("login is required", fasthttp.StatusBadRequest)
	}

	existsLogin, err := a.auth.ExistsUserByLogin(ctx, req.Login)
	if err != nil {
		return nil, err
	}
	if existsLogin {
		return nil, errorx.NewHTTP("login already exists", fasthttp.StatusConflict)
	}

	if req.Password == "" {
		return nil, errorx.NewHTTP("password is required", fasthttp.StatusBadRequest)
	}

	if err := checkValidLogin(req.Login); err != nil {
		return nil, err
	}

	if err := checkValidPassword(req.Password); err != nil {
		return nil, err
	}

	return a.auth.RegisterUser(ctx, req)
}
