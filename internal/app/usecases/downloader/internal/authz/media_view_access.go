package authz

import (
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// HasDownloadsViewAll returns true if the roles allow viewing all downloads (global access)
func (a *Authorization) HasDownloadsViewAll(roles dtypes.UserRoleIDs) bool {
	return a.HasViewAllAccess(roles)
}

// ShouldRestrictDownloads returns true if access to downloads should be restricted by user
func (a *Authorization) ShouldRestrictDownloads(roles dtypes.UserRoleIDs) bool {
	return !a.HasDownloadsViewAll(roles)
}

func (a *Authorization) HasMediaViewAccess(authCtx dauth.AuthContext, media *ddownload.MediaDownload) bool {
	if media == nil {
		return false
	}

	if a.HasViewAllAccess(authCtx.RoleIDs) {
		return true
	}

	if media.Visibility == dtypes.MediaVisibilityPublic {
		return true
	}

	if media.Visibility == dtypes.MediaVisibilityAuthenticated && authCtx.UserType() > dtypes.UserTypeGuest {
		return true
	}

	if media.UserID != nil && authCtx.UserID == *media.UserID {
		return true
	}

	return false
}
