package ddownload

import (
	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type DownloadState struct {
	TaskID *uuid.UUID

	Download *MediaDownload
	Progress *dtypes.DownloadProgress
}

// InitFromMediaDownload initializes the state from a file.
func (s *DownloadState) InitFromMediaDownload(download *MediaDownload) {
	if s == nil {
		return
	}

	if download == nil {
		return
	}

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

	copy.Download = src.Download.Clone()
	copy.Progress = src.Progress.Clone()

	return &copy

}

func (s *DownloadState) IsEmpty() bool {
	if s == nil {
		return true
	}
	return s.Download == nil
}

func (src *DownloadState) DownloadStatus() dtypes.MediaDownloadStatus {
	if src == nil || src.Download == nil {
		return dtypes.MediaDownloadStatusNone
	}
	return src.Download.Status
}

func (src *DownloadState) UserID() uuid.UUID {
	if src == nil {
		return uuid.Nil()
	}
	return src.Download.EnsureUserID()
}

func (src *DownloadState) DownloadID() uuid.UUID {
	if src == nil || src.Download == nil {
		return uuid.Nil()
	}
	return src.Download.DownloadID
}
