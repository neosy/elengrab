package ytdlpsrv

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	dyoutubeinfo "github.com/neosy/elengrab/internal/domain/youtube_info"
)

// GetBestFormat retrieves and parses video best format for the given URL.
func (srv *YtDlpService) GetBestFormat(url string) (*dyoutubeinfo.YouTubeInfo, error) {
	return srv.getBestFormat(url, "b")
}

func (srv *YtDlpService) getBestFormat(url string, format string) (*dyoutubeinfo.YouTubeInfo, error) {
	cmd := exec.Command(srv.cmdPath, "--no-playlist", "-f", format, "--get-format", url)

	// Buffers to capture stdout and stderr
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errOut := fmt.Errorf("%s failed: %v, stderr: %s", ytDlpName, err, stderr.String())
		srv.logger.Error(errOut.Error())
		return nil, errOut
	}

	outStr := strings.TrimSpace(out.String())
	if outStr == "" {
		err := fmt.Errorf("best format not found")
		srv.logger.Error(err.Error())
		return nil, err
	}

	parts := strings.SplitN(outStr, " - ", 2)
	bestFormatId := parts[0]

	info, err := srv.GetFormats(url)
	if err != nil {
		srv.logger.Error(err.Error())
		return nil, err
	}

	var bestFormat *dyoutubeinfo.Format
	for _, f := range info.Formats {
		if f.FormatId == bestFormatId {
			bestFormat = &f
			break
		}
	}

	if bestFormat == nil {
		err := fmt.Errorf("best format not found")
		srv.logger.Error(err.Error())
		return nil, err
	}

	bestInfo := *info
	bestInfo.Formats = []dyoutubeinfo.Format{*bestFormat}

	return &bestInfo, nil
}
