package ucdownloader

import "path/filepath"

func (uc *YouTubeDownloader) GetFilePath(fileName string) string {
	return filepath.Join(uc.downloadsDir, fileName)
}
