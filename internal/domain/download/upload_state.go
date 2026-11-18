package ddownload

import (
	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type UploadState struct {
	Download *MediaDownload
}

func NewUploadState(download *MediaDownload) *UploadState {
	state := &UploadState{}
	state.InitFromMediaDownload(download)

	return state
}

// InitFromMediaDownload initializes the state from a file.
func (s *UploadState) InitFromMediaDownload(download *MediaDownload) {
	if s == nil {
		return
	}

	s.Download = download
}

// Copy creates a deep copy of the UploadState.
func (src *UploadState) Clone() *UploadState {
	if src == nil {
		return nil
	}

	copy := *src

	copy.Download = src.Download.Clone()

	return &copy

}

func (s *UploadState) IsEmpty() bool {
	if s == nil {
		return true
	}
	return s.Download == nil
}

func (src *UploadState) UserID() uuid.UUID {
	if src == nil {
		return uuid.Nil()
	}
	return src.Download.EnsureUserID()
}

func (src *UploadState) DownloadID() uuid.UUID {
	if src == nil || src.Download == nil {
		return uuid.Nil()
	}
	return src.Download.DownloadID
}

func (src *UploadState) Status() dtypes.MediaDownloadStatus {
	if src == nil || src.Download == nil {
		return dtypes.MediaDownloadStatusNone
	}
	return src.Download.Status
}
