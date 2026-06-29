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

func buildAdminUserPath(path string, userID uuid.UUID) string {
	id := idcodec.EncodeUUIDBase64URL(userID)
	return AdminGroup + strings.Replace(path, "{userId}", id, 1)
}

func BuildAdminUserRolesPath(userID uuid.UUID) string {
	return buildAdminUserPath(AdminUserRolesPath, userID)
}

func BuildAdminUserTableRowPath(userID uuid.UUID) string {
	return buildAdminUserPath(AdminUserTableRowPath, userID)
}
