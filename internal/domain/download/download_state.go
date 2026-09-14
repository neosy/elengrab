package ddownload

import (
	"github.com/google/uuid"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type DownloadState struct {
	UserID     *uuid.UUID
	DownloadID uuid.UUID
	TaskID     *uuid.UUID

	Download *MediaDownload
	Progress *dservices.DownloaderProgress
}

// InitFromMediaDownload initializes the state from a file.
func (s *DownloadState) InitFromMediaDownload(download *MediaDownload) {
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

// Copy creates a deep copy of the DownloadState.
func (src *DownloadState) Clone() *DownloadState {
	if src == nil {
		return nil
	}

	copy := *src

	copy.UserID = uptr.Clone(src.UserID)
	copy.TaskID = uptr.Clone(src.TaskID)
	copy.Download = src.Download.Clone()
	copy.Progress = src.Progress.Clone()

	return &copy

}
