package httppaths

import (
	"strings"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
)

const (
	// Groups
	AdminGroup = "/admin"

	// Paths Users
	UsersPath        = "/users"
	UserDetailPath   = UsersPath + "/{userId}/detail"
	UserTableRowPath = UsersPath + "/{userId}/table-row"
	UserRolesPath    = UsersPath + "/{userId}/roles"

	// Paths Groups, Logs, Settings
	UserGroupsPath = "/groups"
	LogsPath       = "/logs"
	SettingsPath   = "/settings"

	// Admin paths
	AdminUsersPath      = AdminGroup + UsersPath
	AdminUserGroupsPath = AdminGroup + UserGroupsPath
)

func BuildAdminUserDetailTemplatePath() string {
	return AdminGroup + UserDetailPath
}

func buildAdminUserPath(path string, userID uuid.UUID) string {
	id := idcodec.EncodeUUIDBase64URL(userID)
	return AdminGroup + strings.Replace(path, "{userId}", id, 1)
}

func BuildAdminUserRolesPath(userID uuid.UUID) string {
	return buildAdminUserPath(UserRolesPath, userID)
}

func BuildAdminUserTableRowPath(userID uuid.UUID) string {
	return buildAdminUserPath(UserTableRowPath, userID)
}
