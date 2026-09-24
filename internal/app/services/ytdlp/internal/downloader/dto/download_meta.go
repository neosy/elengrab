package idto

import (
	"sync"

	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type DownloadMeta struct {
	URL          string
	Title        string
	Description  string
	FileName     string
	FileExt      string
	FileFullName string
	FileSize     *int64

	Channel *dtypes.ChannelSource

	Thumbnail           *dtypes.ImageData
	ThumbnailVideoFrame *dtypes.ImageData

	MediaInfo *dservices.MediaInfo

	Progress *dservices.DownloaderProgress
}

type SafeDownloadMeta struct {
	mu   sync.RWMutex
	Meta *DownloadMeta
}

func (m *SafeDownloadMeta) Lock() {
	m.mu.Lock()
}

func (m *SafeDownloadMeta) Unlock() {
	m.mu.Unlock()
}

func (m *SafeDownloadMeta) RLock() {
	m.mu.RLock()
}

func (m *SafeDownloadMeta) RUnlock() {
	m.mu.RUnlock()
}

func (m *SafeDownloadMeta) CloneMeta() *DownloadMeta {
	m.RLock()
	defer m.RUnlock()

	if m == nil || m.Meta == nil {
		return nil
	}

	metaCopy := *m.Meta
	metaCopy.FileSize = uptr.Clone(m.Meta.FileSize)

	metaCopy.Channel = m.Meta.Channel.Clone()

	metaCopy.Thumbnail = m.Meta.Thumbnail.Clone()
	metaCopy.ThumbnailVideoFrame = m.Meta.ThumbnailVideoFrame.Clone()

	metaCopy.MediaInfo = m.Meta.MediaInfo.Clone()
	metaCopy.Progress = uptr.Clone(m.Meta.Progress)

	return &metaCopy
}

func (m *SafeDownloadMeta) InitialResult() *dservices.DownloaderResult {
	meta := m.CloneMeta()

	var ext = meta.FileExt
	if ext == "" && meta.MediaInfo != nil && meta.MediaInfo.Format != dtypes.FileFormatNone {
		ext = meta.MediaInfo.Format.String()
	}

	return &dservices.DownloaderResult{
		MediaTitle:          meta.Title,
		MediaDescription:    uptr.NonZeroString(meta.Description),
		Filename:            meta.FileName,
		FileExt:             ext,
		FileFullName:        meta.FileFullName,
		Filesize:            meta.FileSize,
		MediaInfo:           meta.MediaInfo,
		Channel:             meta.Channel,
		Thumbnail:           meta.Thumbnail,
		ThumbnailVideoFrame: meta.ThumbnailVideoFrame,
		Progress:            meta.Progress,
	}
}
