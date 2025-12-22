package ddownload

type DownloadResult struct {
	Error error

	ChannelID *string

	YoutubeTitle *string

	// full path to the downloaded file
	FilePath *string

	// file name
	Filename *string

	// ext
	FileExt *string

	// file name
	FileFullName *string

	// file size (byte)
	Filesize *int

	// Fast partial file hash (combined hash of multiple sampled blocks; not a full-file checksum)
	PartialHash *string

	// Channel avatar
	ChannelAvatar *DownloadResultChannelAvatar

	// MediaInfo holds media metadata.
	MediaInfo *MediaInfo
}

type DownloadResultChannelAvatar struct {
	ImageURL    string
	ImageRAW    []byte
	ImageFormat string
}
