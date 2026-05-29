package httppaths

import (
	"strings"

	"github.com/google/uuid"
)

const (
	// Groups
	AdminGroup = "/admin"

	// Paths Users
	AdminUsersPath        = "/users"
	AdminUserDetailPath   = "/users/{userId}/detail"
	AdminUserTableRowPath = "/users/{userId}/table-row"
	AdminUserRolesPath    = "/users/{userId}/roles"

	// Paths Groups, Logs, Settings
	AdminGroupsPath   = "/groups"
	AdminLogsPath     = "/logs"
	AdminSettingsPath = "/settings"
)

func BuildAdminUserDetailTemplatePath() string {
	return AdminGroup + AdminUserDetailPath
}

func BuildAdminUserRolesPath(userID uuid.UUID) string {
	return AdminGroup + strings.Replace(AdminUserRolesPath, "{userId}", userID.String(), 1)
}

func BuildAdminUserTableRowPath(userID uuid.UUID) string {
	return AdminGroup + strings.Replace(AdminUserTableRowPath, "{userId}", userID.String(), 1)
}
