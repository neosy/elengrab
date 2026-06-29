package admin

import (
	"bytes"
	"mime"

	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/admin/dto"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *AdminHandlers) UserDetailHandler(ctx *fasthttp.RequestCtx) {
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

	h.renderUserDetail(ctx, userInfo)
}

func (h *AdminHandlers) renderUserDetail(ctx *fasthttp.RequestCtx, userInfo *ucdto.UserInfoResponse) {
	data, err := h.buildUserDetail(userInfo)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	// Load template
	tmpl, err := h.templates.Clone()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errTemplateInternal(err))
		return
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, components.AmdinUserDetailPanelKey, *data)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errTemplateInternal(err))
		return
	}

	ctx.SetContentType(mime.TypeByExtension(".html"))
	ctx.SetBody(buf.Bytes())
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *AdminHandlers) buildUserDetail(userInfo *ucdto.UserInfoResponse) (*pages.AdminUserDetail, error) {
	roles := make([]pages.AdminUserRoleWithAssign, 0, len(userInfo.RolesWithAssignment))
	for _, role := range userInfo.RolesWithAssignment {
		roles = append(roles, pages.AdminUserRoleWithAssign{
			RoleID:   role.RoleID,
			Name:     role.Name,
			Assigned: role.Assigned,
		})
	}

	return &pages.AdminUserDetail{
		User:            h.mappers.UserToUserOnPage(&userInfo.User, pages.NewAdminUserIcons()),
		RolesWithAssign: roles,
		UpdateRolesPath: httppaths.BuildAdminUserRolesPath(userInfo.User.UserID),
	}, nil
}
