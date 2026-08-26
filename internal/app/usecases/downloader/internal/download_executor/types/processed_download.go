package types

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type ProcessedDownload struct {
	MediaTitle         string
	MediaTitleOriginal string

	MediaDescription         *string
	MediaDescriptionOriginal *string

	Filename     string
	FileFullName string

	FileExt  string
	Filesize *int64

	PartialHash *string

	ChannelID *string

	MediaInfo *dtypes.MediaInfo
}
