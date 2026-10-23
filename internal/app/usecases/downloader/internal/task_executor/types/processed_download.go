package types

import (
	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type ProcessedDownload struct {
	MediaTitle         string
	MediaTitleOriginal string

	MediaDescription         *string
	MediaDescriptionOriginal *string

	FileNameNoExt string
	FileName      string

	FileExtNoDot string
	Filesize     *int64

	PartialHash *string

	ChannelID uuid.UUID

	MediaInfo *dtypes.MediaInfo
}
