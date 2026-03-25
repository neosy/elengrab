package authz

import dtypes "github.com/neosy/elengrab/internal/domain/types"

// AllowFilesViewAll returns true if the roles allow viewing all files (global access)
func (a *Authorization) AllowFilesViewAll(roles []dtypes.UserRole) bool {
	if a.historyMode == dtypes.HistoryModeGlobal {
		return true
	}
	return dtypes.HasAnyRole(roles, dtypes.UserRoleAdmin, dtypes.UserRoleViewerAll)
}

// RestrictFilesByUser returns true if access to files should be restricted by user
func (a *Authorization) RestrictFilesByUser(roles []dtypes.UserRole) bool {
	return !a.AllowFilesViewAll(roles)
}
