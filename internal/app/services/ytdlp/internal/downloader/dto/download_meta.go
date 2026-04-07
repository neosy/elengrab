package idto

import (
	"sync"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type DownloadMeta struct {
	URL          string
	Title        string
	FileName     string
	FileExt      string
	FileFullName string
	FilePath     string
	FileSize     *int64
	ChannelID    *string
	ChannelURL   string
	ChannelTitle string
	MediaInfo    *ddownload.MediaInfo
	Channel      *ddownload.DownloadChannel
	Progress     *ddownload.DownloadProgress
	Options      DownloadOptions
}

type DownloadOptions struct {
	ConcurrentFragments    uint8
	RequiresYouTubeCookies bool
	Extractor              string
	ExtractorArgs          *string
	Args                   []string
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

	metaCopy := *m.Meta
	metaCopy.FileSize = uptr.Copy(m.Meta.FileSize)
	metaCopy.ChannelID = uptr.Copy(m.Meta.ChannelID)
	metaCopy.MediaInfo = m.Meta.MediaInfo.Copy()
	metaCopy.Channel = m.Meta.Channel.Copy()
	metaCopy.Progress = uptr.Copy(m.Meta.Progress)
	metaCopy.Options.ExtractorArgs = uptr.Copy(m.Meta.Options.ExtractorArgs)
	metaCopy.Options.Args = m.Meta.Options.Args[0:len(m.Meta.Options.Args):len(m.Meta.Options.Args)]

	return &metaCopy
}

func (m *SafeDownloadMeta) InitialResult() *ddownload.DownloadResult {
	meta := m.CopyMeta()

	var ext = meta.FileExt
	if ext == "" && meta.MediaInfo != nil && meta.MediaInfo.Format != dtypes.FileFormatNone {
		ext = meta.MediaInfo.Format.String()
	}

	return &ddownload.DownloadResult{
		ChannelID:    meta.ChannelID,
		MediaTitle:   meta.Title,
		FilePath:     meta.FilePath,
		Filename:     meta.FileName,
		FileExt:      ext,
		FileFullName: meta.FileFullName,
		Filesize:     meta.FileSize,
		MediaInfo:    meta.MediaInfo,
		Channel:      meta.Channel,
		Progress:     meta.Progress,
	}
}
