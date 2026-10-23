package ddownload

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"uuid"

	"github.com/neosy/elengrab/internal/app/utils/filename"
	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	dconsts "github.com/neosy/elengrab/internal/domain/consts"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/stringx"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

const (
	// removeTrailingTitleHashtagsEnabled controls whether trailing hashtags are removed from media titles.
	removeTrailingTitleHashtagsEnabled = false

	// validateMediaTitleControlCharactersEnabled controls whether control characters in media titles are rejected.
	validateMediaTitleControlCharactersEnabled = false
)

type MediaDownload struct {
	// Unique file identifier (UUID)
	DownloadID uuid.UUID

	// Unique short download code
	DownloadCode string

	// Source of the media record: a remote URL or a file uploaded by the user
	SourceType dtypes.MediaDownloadSourceType

	// Associated user identifier (UUID)
	UserID *uuid.UUID

	// Status
	Status dtypes.MediaDownloadStatus

	// Media URL
	MediaURL *string

	// Media title
	MediaTitle string

	// Original title from the media source
	MediaTitleOriginal string

	// Media description
	MediaDescription *string

	// Original description from the media source
	MediaDescriptionOriginal *string

	// Channel ID
	ChannelID uuid.UUID

	// Original file name
	FileNameNoExt string

	// File extension without a leading dot (e.g., "mp3", "mp4")
	FileExtNoDot string

	// Full file name including extension
	FileName string

	// File size (byte)
	FileSize *int64

	// Fast partial file hash (combined hash of multiple sampled blocks; not a full-file checksum)
	PartialHash *string

	// Human-readable safe full name
	SafeReadableFileName string

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
	DownloadTask   *DownloadTask
	ProcessingTask *MediaProcessingTask
}

func (m *MediaDownload) IsYouTube() bool {
	if m == nil {
		return false
	}

	url := m.SourceURL()

	if url == "" {
		return false
	}

	return hostdetect.YouTube(url)
}

func (m *MediaDownload) IsVideo() bool {
	if m == nil || m.MediaInfo == nil {
		return false
	}

	return m.MediaInfo.IsVideo()
}

func (src *MediaDownload) Clone() *MediaDownload {
	if src == nil {
		return nil
	}

	copy := *src

	copy.UserID = uptr.Clone(src.UserID)
	copy.MediaURL = uptr.Clone(src.MediaURL)
	copy.MediaDescription = uptr.Clone(src.MediaDescription)
	copy.MediaDescriptionOriginal = uptr.Clone(src.MediaDescriptionOriginal)
	copy.FileSize = uptr.Clone(src.FileSize)
	copy.PartialHash = uptr.Clone(src.PartialHash)
	copy.MediaInfo = src.MediaInfo.Clone()
	copy.ErrorMessage = uptr.Clone(src.ErrorMessage)
	copy.DownloadedAt = uptr.Clone(src.DownloadedAt)
	copy.DeletedAt = uptr.Clone(src.DeletedAt)
	copy.DownloadTask = src.DownloadTask.Clone()

	return &copy
}

// FileExt returns the extension with a leading dot (e.g. ".mp4").
func (d *MediaDownload) FileExt() string {
	if d.FileExtNoDot == "" {
		return ""
	}

	return "." + d.FileExtNoDot
}

// SourceURL returns the media source URL, or an empty string for entries without one (e.g. uploads).
func (d *MediaDownload) SourceURL() string {
	if d.MediaURL == nil {
		return ""
	}
	return *d.MediaURL
}

func (d *MediaDownload) Normalize() {
	d.MediaTitle = strings.TrimSpace(d.MediaTitle)
	d.MediaTitleOriginal = strings.TrimSpace(d.MediaTitleOriginal)

	if removeTrailingTitleHashtagsEnabled {
		d.MediaTitle = stringx.RemoveTrailingHashtags(d.MediaTitle)
		d.MediaTitleOriginal = stringx.RemoveTrailingHashtags(d.MediaTitleOriginal)
	}

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

func (d *MediaDownload) NormalizeFileName() string {
	name := strings.TrimSpace(d.MediaTitle)
	name = fasthttpx.SanitizeFileName(name)
	name = stringx.TruncateBytesWords(name, filename.MaxTitleLength)
	return name + d.FileExt()
}

func (d *MediaDownload) NormalizeForSave() {
	d.Normalize()

	d.MediaTitle = stringx.TruncateWords(d.MediaTitle, dconsts.MediaTitleMaxLength)
	d.MediaTitleOriginal = stringx.TruncateWords(d.MediaTitleOriginal, dconsts.MediaTitleMaxLength)

	if d.MediaDescription != nil {
		*d.MediaDescription = stringx.Truncate(*d.MediaDescription, dconsts.MediaDescriptionMaxLength)
	}

	if d.MediaDescriptionOriginal != nil {
		*d.MediaDescriptionOriginal = stringx.Truncate(*d.MediaDescriptionOriginal, dconsts.MediaDescriptionMaxLength)
	}

	if (d.UserID == nil || *d.UserID == uuid.Nil()) && d.Visibility == dtypes.MediaVisibilityPrivate {
		d.Visibility = dtypes.MediaVisibilityPublic
	}
}

func (d *MediaDownload) Validate() error {
	if d.DownloadCode == "" {
		return errors.New("download code is required")
	}

	if err := d.SourceType.Validate(); err != nil {
		return err
	}

	if d.SourceType == dtypes.MediaDownloadSourceTypeURL {
		if d.SourceURL() == "" {
			return errors.New("media url is required for url source type")
		}
	}

	if d.SourceType == dtypes.MediaDownloadSourceTypeUpload {
		if d.MediaInfo == nil {
			return errors.New("media info is required for upload source type")
		}
	}

	if utf8.RuneCountInString(d.MediaTitle) > dconsts.MediaTitleMaxLength {
		return errorx.NewHTTPMessage(
			fmt.Sprintf("Title must not exceed %d characters", dconsts.MediaTitleMaxLength),
			http.StatusBadRequest,
		)
	}

	// Reject media titles containing control characters when validation is enabled.
	if validateMediaTitleControlCharactersEnabled && strings.IndexFunc(d.MediaTitle, unicode.IsControl) >= 0 {
		return errorx.NewHTTPMessage("Title contains invalid characters", http.StatusBadRequest)
	}

	if d.MediaDescription != nil {
		if utf8.RuneCountInString(*d.MediaDescription) > dconsts.MediaDescriptionMaxLength {
			return errorx.NewHTTPMessage(
				fmt.Sprintf("Description must not exceed %d characters", dconsts.MediaDescriptionMaxLength),
				http.StatusBadRequest,
			)
		}
	}

	return nil
}
