package pages

import (
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/icons"
)

type (
	AdminContentData struct {
		Title string
		Icons *AdminIcons
	}

	AdminContentDashboardData struct {
		AdminContentData
	}

	AdminContentUsersData struct {
		AdminContentData

		UserDetailQueryPath string
		Users               []AdminUser
	}

	AdminIcons struct {
		Users template.HTML
	}

	AdminUser struct {
		UserID string
		Login  string
		Name   string
		Email  string
		Status string
		Roles  string

		LoadTableRowPath string
		Icons            *AdminUserIcons
	}

	AdminUserIcons struct {
		UserEdit template.HTML
	}

	AdminUserDetail struct {
		User            AdminUser
		RolesWithAssign []AdminUserRoleWithAssign

		UpdateRolesPath string
	}

	AdminUserRoleWithAssign struct {
		RoleID   string
		Name     string
		Assigned bool
	}
)

func NewAdminIcons() *AdminIcons {
	return &AdminIcons{
		Users: icons.AdminUsersIcon.FileRaw(),
	}
}

func NewAdminUserIcons() *AdminUserIcons {
	return &AdminUserIcons{
		UserEdit: icons.AdminUserEditIcon.FileRaw(),
	}
}

func NewAdminContentData(title string) AdminContentData {
	return AdminContentData{
		Title: title,
		Icons: NewAdminIcons(),
	}
}
