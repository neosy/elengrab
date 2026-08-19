package ddownload

import (
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/stringx"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

const (
	MediaTitleMaxLength       = 100
	MediaDescriptionMaxLength = 5000
)

type MediaDownload struct {
	// Unique file identifier (UUID)
	DownloadID uuid.UUID

	// Associated user identifier (UUID)
	UserID *uuid.UUID

	// Status
	Status dtypes.MediaDownloadStatus

	// Media URL
	MediaURL string

	// Original title from the media source
	MediaTitleOriginal string

	// Media title
	MediaTitle string

	// Original description from the media source
	MediaDescriptionOriginal *string

	// Media description
	MediaDescription *string

	// Channel ID
	ChannelID *string

	// Original file name
	FileName string

	// File extension (e.g., "mp3", "mp4")
	Ext string

	// Full file name including extension
	FileFullName string

	// File size (byte)
	FileSize *int64

	// Fast partial file hash (combined hash of multiple sampled blocks; not a full-file checksum)
	PartialHash *string

	// Human-readable safe full name
	SafeReadableFileFullName string

	// MediaInfo holds media metadata.
	MediaInfo *dtypes.MediaInfo

	// Error message
	ErrorMessage *string

	// Visibility access level for media (public or private)
	Visibility dtypes.MediaVisibility

	// Downloaded timestamp
	DownloadedAt *time.Time

	// Timestamp when the record was created
	CreatedAt time.Time

	// Timestamp when the record was last updated
	UpdatedAt time.Time

	// Timestamp when the record was soft deleted
	DeletedAt *time.Time

	// Related task
	DownloadTask *DownloadTask
}

func (f *MediaDownload) IsYouTube() bool {
	return hostdetect.YouTube(f.MediaURL)
}

func (src *MediaDownload) Copy() *MediaDownload {
	if src == nil {
		return nil
	}

	copy := uptr.Copy(src)
	copy.UserID = uptr.Copy(src.UserID)
	copy.MediaDescription = uptr.Copy(src.MediaDescription)
	copy.ChannelID = uptr.Copy(src.ChannelID)
	copy.FileSize = uptr.Copy(src.FileSize)
	copy.PartialHash = uptr.Copy(src.PartialHash)
	copy.MediaInfo = src.MediaInfo.Copy()
	copy.ErrorMessage = uptr.Copy(src.ErrorMessage)
	copy.DownloadedAt = uptr.Copy(src.DownloadedAt)
	copy.DeletedAt = uptr.Copy(src.DeletedAt)
	copy.DownloadTask = src.DownloadTask.Copy()

	return copy
}

func (d *MediaDownload) Normalize() {
	d.MediaTitle = strings.TrimSpace(d.MediaTitle)
	d.MediaTitleOriginal = strings.TrimSpace(d.MediaTitleOriginal)

	// TODO: Normalize title hashtags after deciding how to handle platform-specific tags.
	// d.MediaTitle = stringx.RemoveTrailingHashtags(d.MediaTitle)
	// d.MediaTitleOriginal = stringx.RemoveTrailingHashtags(d.MediaTitleOriginal)

	if d.MediaDescription != nil {
		*d.MediaDescription = strings.TrimSpace(*d.MediaDescription)
		if *d.MediaDescription == "" {
			d.MediaDescription = nil
		}
	}

	if d.MediaDescriptionOriginal != nil {
		*d.MediaDescriptionOriginal = strings.TrimSpace(*d.MediaDescriptionOriginal)
		if *d.MediaDescriptionOriginal == "" {
			d.MediaDescriptionOriginal = nil
		}
	}
}

func (d *MediaDownload) NormalizeFileFullName() string {
	name := strings.TrimSpace(d.MediaTitle)
	name = fasthttpx.SanitizeFileName(name)
	name = stringx.TruncateBytesWords(name, MediaTitleMaxLength)
	return name + "." + d.Ext
}

func (d *MediaDownload) NormalizeForSave() {
	d.Normalize()

	d.MediaTitle = stringx.TruncateWords(d.MediaTitle, MediaTitleMaxLength)
	d.MediaTitleOriginal = stringx.TruncateWords(d.MediaTitleOriginal, MediaTitleMaxLength)

	if d.MediaDescription != nil {
		*d.MediaDescription = stringx.Truncate(*d.MediaDescription, MediaDescriptionMaxLength)
	}
	if d.MediaDescriptionOriginal != nil {
		*d.MediaDescriptionOriginal = stringx.Truncate(*d.MediaDescriptionOriginal, MediaDescriptionMaxLength)
	}
}

func (d *MediaDownload) Validate() error {
	if utf8.RuneCountInString(d.MediaTitle) > MediaTitleMaxLength {
		return errorx.NewHTTPMessage("Title must not exceed 100 characters", http.StatusBadRequest)
	}

	// TODO: Determine the proper validation rules for media title.
	// 		 Current control character check may be too restrictive.
	// if strings.IndexFunc(d.MediaTitle, unicode.IsControl) >= 0 {
	// 	return errorx.NewHTTPMessage("Title contains invalid characters", http.StatusBadRequest)
	// }

	if d.MediaDescription != nil {
		description := *d.MediaDescription

		if utf8.RuneCountInString(description) > MediaDescriptionMaxLength {
			return errorx.NewHTTPMessage(
				"Description must not exceed 5000 characters",
				http.StatusBadRequest,
			)
		}

	}

	return nil
}
