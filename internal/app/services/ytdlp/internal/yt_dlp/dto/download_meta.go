package idto

import ddownload "github.com/neosy/elengrab/internal/domain/download"

type DownloadMeta struct {
	Title        string
	FileName     string
	FileExt      string
	FileFullName string
	FilePath     string
	FileSize     *int
	ChannelID    *string
	ChannelURL   string
	MediaInfo    *ddownload.MediaInfo
	Args         []string
}

func (m *DownloadMeta) InitialResult() *ddownload.DownloadResult {
	return &ddownload.DownloadResult{
		ChannelID:    m.ChannelID,
		YoutubeTitle: m.Title,
		FilePath:     m.FilePath,
		Filename:     m.FileName,
		FileExt:      m.FileExt,
		FileFullName: m.FileFullName,
		Filesize:     m.FileSize,
		MediaInfo:    m.MediaInfo,
	}
}
