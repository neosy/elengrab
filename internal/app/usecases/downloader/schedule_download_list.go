package downloader

import (
	"context"
	"uuid"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/neosy/elengrab/internal/app/usecases/types"
	iconfig "github.com/neosy/elengrab/internal/config"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	ierrors "github.com/neosy/elengrab/internal/errors"
)

func (uc *downloader) ScheduleDownloadList(
	ctx context.Context,
	authCtx dauth.AuthContext,
	req *dto.ScheduleDownloadListRequest,
) (*dto.ScheduleDownloadListResponse, error) {
	if req == nil {
		return nil, ierrors.ErrFuncParamNullPointer
	}

	if len(req.Items) == 0 {
		return nil, ierrors.ErrFuncContainsEmptyFields
	}

	if err := uc.validateWriteAccess(authCtx, types.BroadcastNotificationModuleGrabForm); err != nil {
		return nil, err
	}

	broadcastError := uc.newErrorBroadcaster(authCtx, types.BroadcastNotificationModuleGrabForm)

	mediaVisibility := iconfig.InitialMediaVisibility()
	if authCtx.UserID == uuid.Nil() {
		mediaVisibility = dtypes.MediaVisibilityPublic
	}

	createdDownloads := make([]*ddownload.MediaDownload, 0, len(req.Items))

	err := uc.download.Tx(ctx, func(ctx context.Context) error {
		for _, item := range req.Items {
			downloadID := uuid.New()
			uniqueFileNameNoExt := downloadID.String()

			options := req.Options
			options.UniqueFileNameNoExt = &uniqueFileNameNoExt

			err := uc.download.CreateDownload(
				ctx,
				&ddownload.MediaDownload{
					DownloadID:    downloadID,
					SourceType:    dtypes.MediaDownloadSourceTypeURL,
					UserID:        &authCtx.UserID,
					FileNameNoExt: uniqueFileNameNoExt,
					MediaURL:      &item.URL,
					Visibility:    mediaVisibility,
				},
				&options,
			)
			if err != nil {
				uc.logger.Error("Insert record failed", "error", err)
				return err
			}

			download, err := uc.enqueueDownloadByID(ctx, downloadID)
			if err != nil {
				return err
			}

			createdDownloads = append(createdDownloads, download)
		}
		return nil
	})
	if err != nil {
		return nil, broadcastError(err)
	}

	resp := &dto.ScheduleDownloadListResponse{
		Items: make([]dto.ScheduleDownloadListResponseItem, 0, len(createdDownloads)),
	}

	for _, download := range createdDownloads {
		downloadInfo, err := uc.resolveActualDownloadInfoByDownload(ctx, download)
		if err != nil {
			continue
		}

		uc.broadcastMediaDownloadAdded(downloadInfo)

		resp.Items = append(resp.Items,
			dto.ScheduleDownloadListResponseItem{
				DownloadID: download.DownloadID,
				URL:        download.SourceURL(),
				MediaTitle: download.MediaTitle,
			},
		)
	}

	return resp, nil
}
