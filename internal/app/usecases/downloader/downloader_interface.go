package downloader

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/broadcaster"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type downloaderBase interface {
	AppMode() dtypes.AppMode
	DemoMode() bool
}

type Downloader interface {
	downloaderBase

	Start(ctx context.Context)
}

type DownloaderMaintenance interface {
	downloaderBase

	UpdateHash(ctx context.Context) error
	DeleteDuplicates(ctx context.Context) error
	DeleteMissingDownloads(ctx context.Context, enableMoveUnmatchedFiles bool) error
	DeleteFailedDownloads(ctx context.Context) error
	ResetStuckJobs(ctx context.Context) error
}

type DownloaderTask interface {
	downloaderBase

	ExecuteDownloadTask(
		ctx context.Context,
		workerID uint64,
		task *ddownload.DownloadTask,
	) error

	ExecuteRefreshMetadataTask(
		ctx context.Context,
		workerID uint64,
		task *ddownload.RefreshMetadataTask,
	) error

	UpdateSystemInfo()
}

type DownloaderAPI interface {
	downloaderBase

	MigrateGuestData(ctx context.Context, guestID, userID uuid.UUID) error
	Broadcaster() *broadcaster.Broadcaster
	NotifyDownloadUpdated(ctx context.Context, downloadID uuid.UUID)
	NotifyDownloadChanged(ctx context.Context, req dto.MediaDownloadChanged)
	FindYoutubeChannelInfo(ctx context.Context, channelID string) (*dmedia.YoutubeChannel, error)
	GetYoutubeChannelInfo(ctx context.Context, channelID string) (*dmedia.YoutubeChannel, error)
	GetDownloadInfo(ctx context.Context, authCtx dauth.AuthContext, downloadID uuid.UUID) (*dto.MediaDownloadInfo, error)
	GetDownloadInfoUnrestricted(ctx context.Context, downloadID uuid.UUID) (*dto.MediaDownloadInfo, error)
	GetDownloadInfoForEdit(
		ctx context.Context,
		authCtx dauth.AuthContext,
		downloadID uuid.UUID,
	) (*dto.MediaDownloadInfo, error)
	ListDownloadInfo(
		ctx context.Context,
		authCtx dauth.AuthContext,
		query dto.MediaDownloadQuery,
	) ([]*dto.MediaDownloadInfo, error)
	PatchMediaDownload(ctx context.Context, authCtx dauth.AuthContext, req dto.PatchMediaDownloadRequest) error
	DeleteDownload(ctx context.Context, authCtx dauth.AuthContext, downloadID uuid.UUID) error
	HasWriteOperation(authCtx dauth.AuthContext) bool
	CanAddMediaDownload(authCtx dauth.AuthContext) bool
	GetLastWatchPosition(ctx context.Context, authCtx dauth.AuthContext, downloadID uuid.UUID) (time.Duration, error)
	ScheduleDownload(
		ctx context.Context,
		authCtx dauth.AuthContext,
		url string,
		options *ddownload.DownloadOptions,
	) (*dto.ScheduleDownloadResponse, error)
	RetryDownload(
		ctx context.Context,
		authCtx dauth.AuthContext,
		downloadID uuid.UUID,
	) (*dto.MediaDownloadInfo, error)
	ScheduleRefreshMetadata(
		ctx context.Context,
		authCtx dauth.AuthContext,
		downloadID uuid.UUID,
	) error
	GetDownloadImage(
		ctx context.Context,
		userCtx dauth.AuthContext,
		downloadID uuid.UUID,
		sources []dtypes.ImageSource,
	) (*dtypes.ImageData, error)
	TrackMediaWatchEvent(ctx context.Context, authCtx dauth.AuthContext, req dto.TrackMediaWatchEventRequest) error
	SystemInfo() dto.SystemInfoResponse
}

type InternalDownloader interface {
	downloaderBase

	// Usecases
	MediaDownload() MediaDownload
	MediaWatch() MediaWatch
}
