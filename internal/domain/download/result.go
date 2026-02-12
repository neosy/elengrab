package ddownload

type DownloadResult struct {
	Error error

	ChannelID *string

	MediaTitle string

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
	Channel *DownloadChannel

	// MediaInfo holds media metadata.
	MediaInfo *MediaInfo

	// Download progress
	Progress *DownloadProgress
}
