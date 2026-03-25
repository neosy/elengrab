package ddownload

import (
	"github.com/google/uuid"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type DownloadState struct {
	UserID *uuid.UUID
	FileID uuid.UUID
	TaskID *uuid.UUID

	File     *File
	Progress *DownloadProgress
}

// InitFromFile initializes the state from a file.
func (s *DownloadState) InitFromFile(file *File) {
	if s == nil {
		return
	}

	if file == nil {
		return
	}

	s.UserID = file.UserID
	s.FileID = file.FileID
	s.File = file

	if file.DownloadTask != nil {
		s.TaskID = &file.DownloadTask.TaskID
	}
}

// InitFromDownloadResult initializes the state from a downloadResult.
func (s *DownloadState) InitFromDownloadResult(result *DownloadResult) {
	if s == nil {
		return
	}

	if result == nil {
		return
	}

	if result.ChannelID != nil && result.Channel != nil {
		s.File.YoutubeChannelID = result.ChannelID
	}

	s.File.MediaTitle = result.MediaTitle
	s.File.Ext = result.FileExt

	if result.Filesize != nil {
		s.File.FileSize = result.Filesize
	}

	if result.PartialHash != nil {
		s.File.PartialHash = result.PartialHash
	}

	if result.MediaInfo != nil {
		s.File.MediaInfo = result.MediaInfo
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
	copy.File = src.File.Copy()
	copy.Progress = src.Progress.Copy()

	return copy

}
