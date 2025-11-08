package ddownload

type DownloadResponse struct {
	Title string
	// full path to the downloaded file
	FilePath string
	// file name
	Filename string
	// Ext
	FileExt string
	// file name
	FileFullName string
}
