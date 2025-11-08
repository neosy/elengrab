package ytdlpsrv

import (
	"fmt"
	"os/exec"
	"path"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

const (
	formatTypeDefault = dtypes.FormatTypeVideoAudio
)

func (srv *YtDlpService) Download(url string, options *ddownload.DownloadOptions) (*ddownload.DownloadResponse, error) {
	var (
		formatType  = formatTypeDefault
		downloadDir = srv.downloadsDir
		fileName    string
		fileExt     string
		filePath    string
		cmd         *exec.Cmd
		args        []string
		title       string
	)

	if options != nil {
		if options.FormatType != dtypes.FormatTypeNone {
			formatType = options.FormatType
		}

		if options.Filename != nil {
			fileName = fileNameWithoutExt(*options.Filename)
		}

		if options.DownloadsDir != nil {
			downloadDir = *options.DownloadsDir
		}
	}

	if err := checkDir(downloadDir); err != nil {
		srv.logger.Error(err.Error())
		return nil, err
	}

	switch formatType {
	case dtypes.FormatTypeVideoAudio, dtypes.FormatTypeVideoOnly:
		info, err := srv.getBestFormat(url, "b")
		if err != nil {
			srv.logger.Error(err.Error())
			return nil, err
		}
		args = append(args, "-f", "bestvideo+bestaudio")
		args = append(args, "--merge-output-format", info.Formats[0].FileExt)
		title = info.Title
		fileExt = info.Formats[0].FileExt
	case dtypes.FormatTypeAudioOnly:
		info, err := srv.getBestFormat(url, "bestaudio")
		if err != nil {
			srv.logger.Error(err.Error())
			return nil, err
		}
		args = append(args, "-f", "bestaudio")
		args = append(args, "--extract-audio", "--audio-format", "mp3")
		args = append(args, "--audio-quality", "0")
		title = info.Title
		fileExt = "mp3"
	}

	if title == "" {
		var err error
		title, err = srv.GetTitle(url)
		if err != nil {
			srv.logger.Error(err.Error())
			return nil, err
		}
	}

	if fileName == "" {
		fileName = uuid.New().String()
	}

	FileFullName := fmt.Sprintf("%s.%s", fileName, fileExt)
	filePath = path.Join(downloadDir, FileFullName)

	// output file
	args = append(args, "-o", filePath)

	// video url
	args = append(args, url)

	cmd = exec.Command(srv.cmdPath, args...)
	outByte, err := cmd.CombinedOutput()
	if err != nil {
		// Log the full output (stdout + stderr)
		output := string(outByte)
		srv.logger.Error("yt-dlp failed", "error", err, "output", output)

		return nil, fmt.Errorf("%s error: %v, output: %s", ytDlpName, err, output)
	}

	srv.logger.Debug(string(outByte))

	resp := &ddownload.DownloadResponse{
		Title:        title,
		FilePath:     filePath,
		Filename:     fileName,
		FileExt:      fileExt,
		FileFullName: FileFullName,
	}

	srv.logger.Info("Download successful", "info", resp)

	return resp, nil
}
