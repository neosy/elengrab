package idto

import (
	"sync"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

type DownloadMeta struct {
	Title         string
	FileName      string
	FileExt       string
	FileFullName  string
	FilePath      string
	FileSize      *int
	ChannelID     *string
	ChannelURL    string
	MediaInfo     *ddownload.MediaInfo
	ChannelAvatar *ddownload.DownloadResultChannelAvatar
	Progress      *ddownload.DownloadProgress
	Args          []string
}

type SafeDownloadMeta struct {
	mu   sync.Mutex
	Meta *DownloadMeta
}

func (m *SafeDownloadMeta) Lock() {
	m.mu.Lock()
}

func (m *SafeDownloadMeta) Unlock() {
	m.mu.Unlock()
}

func (m *SafeDownloadMeta) CopyMeta() *DownloadMeta {
	m.Lock()
	defer m.Unlock()

	if m == nil || m.Meta == nil {
		return nil
	}

	meta := *m.Meta
	meta.FileSize = uptr.Copy(m.Meta.FileSize)
	meta.ChannelID = uptr.Copy(m.Meta.ChannelID)
	meta.MediaInfo = m.Meta.MediaInfo.Copy()
	meta.ChannelAvatar = m.Meta.ChannelAvatar.Copy()
	meta.Progress = uptr.Copy(m.Meta.Progress)
	meta.Args = m.Meta.Args[0:len(m.Meta.Args):len(m.Meta.Args)]

	return &meta
}

func (m *SafeDownloadMeta) InitialResult() *ddownload.DownloadResult {
	meta := m.CopyMeta()
	return &ddownload.DownloadResult{
		ChannelID:     meta.ChannelID,
		YoutubeTitle:  meta.Title,
		FilePath:      meta.FilePath,
		Filename:      meta.FileName,
		FileExt:       meta.FileExt,
		FileFullName:  meta.FileFullName,
		Filesize:      meta.FileSize,
		MediaInfo:     meta.MediaInfo,
		ChannelAvatar: meta.ChannelAvatar,
		Progress:      meta.Progress,
	}
}
