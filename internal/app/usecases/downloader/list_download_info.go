package downloader

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
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
	var options dtypes.QueryMediaOptions

	options.Before = new(query.Before)
	options.Limit = new(query.Limit)

	if uc.authz.ShouldRestrictDownloads(authCtx.RoleIDs) {
		options.Filters.Add(dtypes.QueryFilterNameUserID, authCtx.UserID)
		if authCtx.IsUser() {
			options.Visibility = new(dtypes.QueryMediaVisibilityAuthenticated)
		}
	}

	if query.Filters.Title != "" {
		options.Filters.Add(dtypes.QueryFilterNameTitle, query.Filters.Title)
	}

	return uc.listDownloadInfo(ctx, authCtx, options, withAuth(authCtx))
}

func (uc *downloader) listDownloadInfo(
	ctx context.Context,
	authCtx dauth.AuthContext,
	queryOptions dtypes.QueryMediaOptions,
	opts ...callOption,
) ([]*dto.MediaDownloadInfo, error) {
	var (
		mu sync.Mutex
	)

	downloadIDs, err := uc.searchIndex.GetDownloadIDsFromMediaSourceIndex(ctx, &queryOptions)
	if err != nil {
		return nil, err
	}

	if len(downloadIDs) == 0 {
		return nil, nil
	}

	downloads, err := uc.download.GetByIDs(ctx, downloadIDs)
	if err != nil {
		return nil, err
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
				resp, err := uc.findActualDownloadInfoByDownload(ctx, download, opts...)
				if err != nil {
					uc.logger.Warn(
						"Failed to load download info",
						"downloadID", download.DownloadID,
						"error", err,
					)
					return nil
				}

				resp.HasWriteAccess = uc.HasWriteOperation(authCtx)

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
