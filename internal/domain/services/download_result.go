package dservices

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type DownloadResult struct {
	Error error

	ChannelID *string

	MediaTitle       string
	MediaDescription *string

	// full path to the downloaded file
	FilePath string

	// file name
	Filename string

	// ext
	FileExt string

	// file name
	FileFullName string

	// file size (byte)
	Filesize *int64

	// Fast partial file hash (combined hash of multiple sampled blocks; not a full-file checksum)
	PartialHash *string

	// Channel
	Channel *dtypes.Channel

	// MediaInfo holds media metadata.
	MediaInfo *MediaInfo

	// Thumbnail is the thumbnail image for the media, if available.
	Thumbnail *dtypes.ImageData

	// ThumbnailVideoFrame holds the best video frame extracted from the media file, if available.
	ThumbnailVideoFrame *dtypes.ImageData

	// Download progress
	Progress *DownloadProgress
}

func (r *DownloadResult) ProgressChanged(last *DownloadResult) bool {
	if r == nil || last == nil {
		return false
	}

	return (r.Progress != nil && last.Progress == nil && r.Progress.Percent() > 0) ||
		(r.Progress != nil && last.Progress != nil &&
			int(r.Progress.Percent())-int(last.Progress.Percent()) > 1 &&
			int(r.Progress.Percent()) < 100)
}

func (r *DownloadResult) MetadataChanged(last *DownloadResult) bool {
	if r == nil {
		return false
	}

	if last == nil {
		return true
	}

	isChanged := false

	isChanged = isChanged ||
		(r.ChannelID != nil && last.ChannelID == nil) ||
		(r.ChannelID != nil && last.ChannelID != nil && *r.ChannelID != *last.ChannelID) ||
		(r.Channel != nil && last.Channel == nil) ||
		(r.MediaTitle != last.MediaTitle) ||
		(r.MediaDescription != nil && last.MediaDescription == nil) ||
		(r.MediaDescription != nil && last.MediaDescription != nil && *r.MediaDescription != *last.MediaDescription) ||
		(r.Filesize != nil && last.Filesize == nil) ||
		(r.Filesize != nil && last.Filesize != nil && *r.Filesize != *last.Filesize) ||
		(r.MediaInfo != nil && last.MediaInfo == nil) ||
		(r.MediaInfo != nil && last.MediaInfo != nil && r.MediaInfo.Format != last.MediaInfo.Format) ||
		(r.Progress != nil && last.Progress == nil) ||
		(r.Progress != nil && last.Progress != nil &&
			int(r.Progress.Percent()) != int(last.Progress.Percent()) &&
			int(r.Progress.Percent()) >= 100) ||
		(r.Thumbnail != nil && last.Thumbnail == nil) ||
		(r.ThumbnailVideoFrame != nil && last.ThumbnailVideoFrame == nil)

	return isChanged
}
