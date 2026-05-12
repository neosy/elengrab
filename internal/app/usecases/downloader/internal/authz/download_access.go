package authz

import dtypes "github.com/neosy/elengrab/internal/domain/types"

// AllowDownloadsViewAll returns true if the roles allow viewing all downloads (global access)
func (a *Authorization) AllowDownloadsViewAll(roles []dtypes.UserRole) bool {
	if a.appMode == dtypes.AppModePublic {
		return true
	}
	return dtypes.HasAnyRole(roles, dtypes.UserRoleAdmin, dtypes.UserRoleViewerAll)
}

// RestrictDownloadsByUser returns true if access to downloads should be restricted by user
func (a *Authorization) RestrictDownloadsByUser(roles []dtypes.UserRole) bool {
	return !a.AllowDownloadsViewAll(roles)
}
