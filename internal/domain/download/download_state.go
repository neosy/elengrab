package ddownload

import (
	"github.com/google/uuid"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

type DownloadState struct {
	UserID *uuid.UUID
	FileId uuid.UUID
	TaskId *uuid.UUID

	File     *File
	Progress *DownloadProgress
}

func (s *DownloadState) InitFromFile(file *File) {
	if s == nil {
		return
	}

	if file == nil {
		return
	}

	s.UserID = file.UserID
	s.FileId = file.FileId
	s.File = file

	if file.DownloadTask != nil {
		s.TaskId = &file.DownloadTask.TaskId
	}
}

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

func (src *DownloadState) Copy() *DownloadState {
	if src == nil {
		return nil
	}

	copy := uptr.Copy(src)
	copy.UserID = uptr.Copy(src.UserID)
	copy.TaskId = uptr.Copy(src.TaskId)
	copy.File = src.File.Copy()
	copy.Progress = src.Progress.Copy()

	return copy

}
