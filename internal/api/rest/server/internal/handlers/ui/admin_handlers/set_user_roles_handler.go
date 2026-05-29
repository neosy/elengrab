package handlers

import (
	"encoding/json"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin_handlers/dto"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *AdminHandlers) SetUserRolesHandler(ctx *fasthttp.RequestCtx) {
	var postReq = dto.SetUserRolesRequest{}

	err := json.Unmarshal(ctx.PostBody(), &postReq)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	err = h.validators.Validate.Struct(postReq)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	ucReq, err := h.mappers.MapSetUserRolesRequestToUsecase(postReq)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	err = h.usecases.admin.SetUserRoles(ctx, ucReq)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	nfasthttp.WriteResponse(ctx, nil)
}
