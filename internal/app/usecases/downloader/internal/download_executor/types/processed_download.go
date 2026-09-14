package types

import (
	"github.com/google/uuid"
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

	ChannelID uuid.UUID

	MediaInfo *dtypes.MediaInfo
}
