package downloader

import (
	"context"
	"uuid"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (u *downloader) CountUserDownloads(
	ctx context.Context,
	authCtx dauth.AuthContext,
	userID uuid.UUID,
) (int, error) {
	if userID == dauth.AnonymousUserID() {
		return 0, nil
	}

	options := dtypes.NewQueryMediaOptions()
	options.Filters.Add(dtypes.QueryFilterNameUserID, userID)

	if userID == authCtx.UserID || !u.authz.ShouldRestrictDownloads(authCtx.RoleIDs) {
		options.Visibility = new(dtypes.QueryMediaVisibilityAll)
	}

	if options.Visibility == nil {
		if authCtx.IsUserOrHigher() {
			options.Visibility = new(dtypes.QueryMediaVisibilityAuthenticated)
		} else {
			options.Visibility = new(dtypes.QueryMediaVisibilityPublic)
		}
	}

	return u.searchIndex.CountFromIndex(ctx, &options)
}

func (u *downloader) CountChannelDownloads(
	ctx context.Context,
	authCtx dauth.AuthContext,
	channelID uuid.UUID,
) (int, error) {
	if channelID == uuid.Nil() {
		return 0, nil
	}

	options := dtypes.NewQueryMediaOptions()
	options.Filters.Add(dtypes.QueryFilterNameChannelID, channelID)

	if !u.authz.ShouldRestrictDownloads(authCtx.RoleIDs) {
		options.Visibility = new(dtypes.QueryMediaVisibilityAll)
	}

	if options.Visibility == nil {
		if authCtx.IsUserOrHigher() {
			options.Visibility = new(dtypes.QueryMediaVisibilityAuthenticated)
		} else {
			options.Visibility = new(dtypes.QueryMediaVisibilityPublic)
		}
	}

	return u.searchIndex.CountFromIndex(ctx, &options)
}
