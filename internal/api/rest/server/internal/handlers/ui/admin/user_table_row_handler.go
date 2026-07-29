package admin

import (
	"bytes"
	"mime"

	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/admin/dto"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *AdminHandlers) UserTableRowHandler(ctx *fasthttp.RequestCtx) {
	userIDStr, ok := ctx.UserValue(userIDKey).(string)
	if !ok || userIDStr == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrUserIDIsRequired)
		return
	}

	userID, err := idcodec.DecodeUUIDBase64URL(userIDStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrUserIDIsIncorrect.Wrap(err))
		return
	}

	userInfo, err := h.usecases.admin.GetUserInfo(ctx, userID)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	h.renderUserTableRow(ctx, userInfo)
}

func (h *AdminHandlers) renderUserTableRow(ctx *fasthttp.RequestCtx, userInfo *ucdto.UserInfoResponse) {
	data := h.mappers.UserToUserOnPage(&userInfo.User, pages.NewAdminUserIcons())

	var buf bytes.Buffer
	err := h.templates.Base.ExecuteTemplate(&buf, components.AdminUsersPanelUserTableRowKey, data)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errTemplateInternal(err))
		return
	}

	ctx.SetContentType(mime.TypeByExtension(".html"))
	ctx.SetBody(buf.Bytes())
	ctx.SetStatusCode(fasthttp.StatusOK)
}
