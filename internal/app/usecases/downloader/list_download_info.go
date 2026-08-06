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
func (uc *Downloader) ListDownloadInfo(
	ctx context.Context,
	authCtx dauth.AuthContext,
	query dto.MediaDownloadQuery,
) ([]*dto.MediaDownloadInfo, error) {
	options := dtypes.QueryOptions{
		Before: new(query.Before),
		Limit:  new(query.Limit),
	}

	filters := make(map[string]any)
	if uc.authz.ShouldRestrictDownloads(authCtx.RoleIDs) {
		filters[dtypes.QueryFilterNameUserID] = authCtx.UserID
		if authCtx.IsUser() {
			options.MediaVisibility = new(dtypes.QueryMediaVisibilityAuthenticated)
		}
	}

	if query.Filters.Title != "" {
		filters[dtypes.QueryFilterNameTitle] = query.Filters.Title
	}

	return uc.listDownloadInfo(ctx, authCtx, options, filters, withAuth(authCtx))
}

func (uc *Downloader) listDownloadInfo(
	ctx context.Context,
	authCtx dauth.AuthContext,
	queryOptions dtypes.QueryOptions,
	filters map[string]any,
	opts ...callOption,
) ([]*dto.MediaDownloadInfo, error) {
	var (
		mu sync.Mutex
	)

	downloads, err := uc.download.GetBeforeTime(ctx, queryOptions, filters)
	if err != nil {
		uc.logger.Warn("Failed get downloads", "queryOptions", queryOptions, "error", err)
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
