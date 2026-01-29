package downloader

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (d *Downloader) watchProgress(fileSize *int64, std io.Reader, outBuf *bytes.Buffer, onProgressUpdate func(ddownload.DownloadProgress)) {
	var (
		downloadedMap = make(map[int64]int64)
		scanner       = bufio.NewScanner(std)
	)

	for scanner.Scan() {
		line := scanner.Text()

		// Parse progress line
		parts := strings.Split(line, "|")
		if len(parts) != 4 {
			// Output line to buffer
			outBuf.WriteString(line)
			outBuf.WriteByte('\n')
			continue
		}

		// Parse progress parts
		downloaded, _ := strconv.ParseInt(parts[0], 10, 64)
		total, _ := strconv.ParseInt(parts[1], 10, 64)
		eta, _ := strconv.Atoi(parts[2])
		speedFloat, _ := strconv.ParseFloat(parts[3], 64)
		speed := int64(speedFloat)

		// Update downloaded map
		downloadedMap[total] = downloaded
		// Initialize downloaded and total
		downloaded = 0
		if fileSize != nil {
			total = *fileSize
		}

		// Aggregate total and downloaded bytes
		for k, v := range downloadedMap {
			downloaded += v
			if fileSize == nil {
				total += k
			}
		}

		// Call progress update callback
		onProgressUpdate(
			ddownload.DownloadProgress{
				DownloadedBytes:  downloaded,
				TotalBytes:       total,
				ETASeconds:       eta,
				SpeedBytesPerSec: speed,
			})
	}
}
