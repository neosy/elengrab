package ddownload

type DownloadResult struct {
	Error error

	YoutubeTitle string

	// full path to the downloaded file
	FilePath string
	// file name
	Filename string
	// ext
	FileExt string
	// file name
	FileFullName string

	// file size (byte)
	Filesize *int
}
