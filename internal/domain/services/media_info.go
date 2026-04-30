package dservices

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type MediaInfo struct {
	FormatType dtypes.FormatType
	Format     dtypes.FileFormat

	VideoInfo *dtypes.VideoInfo
	AudioInfo *dtypes.AudioInfo
}

func (m MediaInfo) MediaInfoDomain() dtypes.MediaInfo {
	return dtypes.MediaInfo{
		FormatType: m.FormatType,
		Format:     m.Format,
		VideoInfo:  uptr.Copy(m.VideoInfo),
		AudioInfo:  uptr.Copy(m.AudioInfo),
	}
}

func (m *MediaInfo) MediaInfoDomainPtr() *dtypes.MediaInfo {
	if m == nil {
		return nil
	}
	return new(m.MediaInfoDomain())
}

func (m *MediaInfo) Copy() *MediaInfo {
	if m == nil {
		return nil
	}

	info := *m

	info.VideoInfo = uptr.Copy(m.VideoInfo)
	info.AudioInfo = uptr.Copy(m.AudioInfo)

	return &info
}

// HasVideo
func (m *MediaInfo) HasVideo() bool {
	return m.VideoInfo != nil
}
