package ddownload

import (
	"github.com/google/uuid"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type DownloadState struct {
	UserID     *uuid.UUID
	DownloadID uuid.UUID
	TaskID     *uuid.UUID

	Download *MediaDownload
	Progress *dservices.DownloaderProgress
}

// InitFromFile initializes the state from a file.
func (s *DownloadState) InitFromFile(download *MediaDownload) {
	if s == nil {
		return
	}

	if download == nil {
		return
	}

	s.UserID = download.UserID
	s.DownloadID = download.DownloadID
	s.Download = download

	if download.DownloadTask != nil {
		s.TaskID = &download.DownloadTask.TaskID
	}
}

// InitFromDownloaderResult initializes the state from a downloadResult.
func (s *DownloadState) InitFromDownloaderResult(result *dservices.DownloaderResult, mediaInfo *dtypes.MediaInfo) {
	if s == nil {
		return
	}

	if result == nil {
		return
	}

	if result.ChannelID != nil && result.Channel != nil {
		s.Download.ChannelID = result.ChannelID
	}

	s.Download.MediaTitle = result.MediaTitle
	s.Download.MediaDescription = result.MediaDescription
	s.Download.Ext = result.FileExt

	if result.Filesize != nil {
		s.Download.FileSize = result.Filesize
	}

	if result.PartialHash != nil {
		s.Download.PartialHash = result.PartialHash
	}

	if mediaInfo != nil {
		s.Download.MediaInfo = mediaInfo
	}

	if result.Progress != nil {
		s.Progress = result.Progress
	}
}

// Copy creates a deep copy of the DownloadState.
func (src *DownloadState) Copy() *DownloadState {
	if src == nil {
		return nil
	}

	copy := uptr.Copy(src)
	copy.UserID = uptr.Copy(src.UserID)
	copy.TaskID = uptr.Copy(src.TaskID)
	copy.Download = src.Download.Copy()
	copy.Progress = src.Progress.Copy()

	return copy

}
