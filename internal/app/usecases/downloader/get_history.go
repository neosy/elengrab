package downloader

import (
	"context"
	"time"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// LoadHistory retrieves the download history for a user.
func (uc *Downloader) LoadHistory(
	ctx context.Context,
	authCtx dauth.AuthContext,
	before time.Time,
	limit uint64,
	filterTitle string,
) ([]*dto.GetMediaDownloadInfoResponse, error) {
	options := dtypes.QueryOptions{
		Before: new(before),
		Limit:  new(limit),
	}

	filters := make(map[string]any)
	if uc.authz.ShouldRestrictDownloads(authCtx.RoleIDs) {
		filters[dtypes.QueryFilterNameUserID] = authCtx.UserID
		if authCtx.IsUser() {
			options.MediaVisibility = new(dtypes.QueryMediaVisibilityAuthenticated)
		}
	}

	if filterTitle != "" {
		filters[dtypes.QueryFilterNameTitle] = filterTitle
	}

	return uc.getDownloadsInfo(ctx, authCtx, options, filters, withAuth(authCtx))
}
