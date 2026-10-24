package downloader

import (
	"context"
	"sync"

	"uuid"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/workerpool"
)

const listDownloadInfoMaxWorkers = 5

// ListDownloadInfo retrieves the download history for a user.
func (uc *downloader) ListDownloadInfo(
	ctx context.Context,
	authCtx dauth.AuthContext,
	query dto.MediaDownloadQuery,
) ([]*dto.MediaDownloadInfo, error) {
	options := dtypes.NewQueryMediaOptions()

	options.ViewMode = query.ViewMode
	if !options.ViewMode.Exists() {
		options.ViewMode = dtypes.QueryMediaViewModeDefault
	}

	options.LastRecord = query.LastRecord
	options.Limit = new(query.Limit)
	options.IsGuestRequest = authCtx.IsGuest()

	var existsUserIDFilter bool
	if query.Filters != nil {
		if filter, exists := query.Filters.Find(dtypes.QueryFilterNameUserID); exists {
			userID, ok := filter.Value().(uuid.UUID)
			existsUserIDFilter = ok && userID != uuid.Nil()

			if existsUserIDFilter && uc.authz.ShouldRestrictDownloads(authCtx.RoleIDs) && userID != authCtx.UserID {
				options.Visibility = new(dtypes.QueryMediaVisibilityPublic)
			}
		}
	}

	if !existsUserIDFilter && uc.authz.ShouldRestrictDownloads(authCtx.RoleIDs) {
		options.Filters.Add(dtypes.QueryFilterNameUserID, authCtx.UserID)
		if authCtx.IsRegularUser() {
			options.Visibility = new(dtypes.QueryMediaVisibilityAuthenticated)
		}
	}

	for _, filter := range query.Filters.List() {
		options.Filters.Append(filter)
	}

	return uc.listDownloadInfo(ctx, authCtx, options, withAuth(authCtx))
}

func (uc *downloader) listDownloadInfo(
	ctx context.Context,
	authCtx dauth.AuthContext,
	queryOptions dtypes.QueryMediaOptions,
	opts ...callOption,
) ([]*dto.MediaDownloadInfo, error) {
	var mu sync.Mutex

	downloads, err := uc.fetchDownloadsByIndex(ctx, &queryOptions)
	if err != nil {
		return nil, err
	}

	if len(downloads) == 0 {
		return nil, nil
	}

	for i, d := range downloads {
		if d.Status.IsProcessing() {
			state, _ := uc.download.FindState(ctx, d.DownloadID)
			if state != nil {
				downloads[i] = state.Download
			}
		}
	}

	downloadsCount := len(downloads)

	if downloadsCount == 0 {
		return nil, nil
	}

	wp := workerpool.NewWorkerPool(
		nil, "listDownloadInfo",
		workerpool.WithMaxWorkers(min(listDownloadInfoMaxWorkers, uint32(downloadsCount))),
	)
	if err := wp.Start(ctx); err != nil {
		return nil, err
	}
	defer wp.Stop()

	orderedDownloadIDs := make([]uuid.UUID, 0, downloadsCount)
	downloadsInfoByID := make(map[uuid.UUID]*dto.MediaDownloadInfo, downloadsCount)

	for _, download := range downloads {
		orderedDownloadIDs = append(orderedDownloadIDs, download.DownloadID)
	}

	for _, download := range downloads {
		job := workerpool.NewSimpleJob(
			download.DownloadID.String(), "findDownloadInfo",
			func(ctx context.Context, workerID uint64) error {
				resp, err := uc.resolveActualDownloadInfoByDownload(ctx, download, opts...)
				if err != nil {
					uc.logger.Warn(
						"Failed to load download info",
						"downloadID", download.DownloadID,
						"error", err,
					)
					return nil
				}

				resp.HasEditAccess = uc.HasDownloadEditAccess(authCtx, download)
				resp.HasDeleteAccess = uc.HasDownloadDeleteAccess(authCtx, download)

				mu.Lock()
				downloadsInfoByID[download.DownloadID] = resp
				mu.Unlock()

				return nil
			},
		)
		wp.AddJob(job)
	}

	wp.WaitJobs()

	results := make([]*dto.MediaDownloadInfo, 0, downloadsCount)

	for _, id := range orderedDownloadIDs {
		result := downloadsInfoByID[id]
		if result != nil {
			results = append(results, result)
		}
	}

	return results, nil
}

// fetchDownloadsByIndex reads download IDs from the search index and loads
// the matching downloads. The index is the primary source, so entries without
// a download (orphans) are skipped and the index is read again with a larger
// limit until wantLimit downloads are collected or the index is exhausted.
// The caller's options are not modified.
func (uc *downloader) fetchDownloadsByIndex(
	ctx context.Context,
	queryOptions *dtypes.QueryMediaOptions,
) ([]*ddownload.MediaDownload, error) {
	if queryOptions == nil {
		return nil, apperrors.ErrFuncParamNullPointer
	}

	opts := *queryOptions

	var wantLimit uint64
	if opts.Limit != nil {
		wantLimit = *opts.Limit
	}

	var (
		downloads []*ddownload.MediaDownload
		passes    int
		lastIDs   []uuid.UUID
	)

	limit := wantLimit

	for {
		passes++

		downloadIDs, err := uc.searchIndex.GetDownloadIDsFromMediaSourceIndex(ctx, &opts)
		if err != nil {
			return nil, err
		}

		if len(downloadIDs) == 0 {
			break
		}

		lastIDs = downloadIDs

		downloads, err = uc.download.GetByIDs(ctx, downloadIDs)
		if err != nil {
			return nil, err
		}

		// No limit, or the index returned less than asked: nothing more to read.
		if limit == 0 || uint64(len(downloadIDs)) < limit {
			break
		}

		if uint64(len(downloads)) >= wantLimit {
			break
		}

		// Ask the index for as many extra entries as are missing.
		limit += wantLimit - uint64(len(downloads))
		opts.Limit = &limit
	}

	if len(downloads) != len(lastIDs) {
		// Collect the IDs that the index returned but the download storage did not.
		found := make(map[uuid.UUID]struct{}, len(downloads))
		for _, d := range downloads {
			found[d.DownloadID] = struct{}{}
		}

		var orphanIDs []uuid.UUID
		for _, id := range lastIDs {
			if _, ok := found[id]; !ok {
				orphanIDs = append(orphanIDs, id)
			}
		}

		uc.logger.Warn("Orphaned index entries hit during page fetch",
			"wanted", wantLimit,
			"passes", passes,
			"index_ids", len(lastIDs),
			"found", len(downloads),
			"orphans", len(orphanIDs),
			"orphan_ids", orphanIDs,
		)
	}

	if len(downloads) == 0 {
		return nil, nil
	}

	if wantLimit > 0 && uint64(len(downloads)) > wantLimit {
		downloads = downloads[:wantLimit:wantLimit]
	}

	return downloads, nil
}
