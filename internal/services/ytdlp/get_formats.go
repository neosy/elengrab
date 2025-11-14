package ytdlpsrv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	dyoutubeinfo "github.com/neosy/elengrab/internal/domain/youtube_info"
	"github.com/neosy/elengrab/internal/services/ytdlp/dto"
)

// GetFormats retrieves and parses video formats for the given URL.
func (srv *YtDlpService) GetFormats(url string) (*dyoutubeinfo.YouTubeInfo, error) {
	// Prepare command to fetch video info in JSON
	cmd := exec.Command(srv.cmdPath, "--no-playlist", "-J", url)

	// Capture combined stdout and stderr
	out, err := cmd.CombinedOutput()
	if err != nil {
		errOut := fmt.Errorf("%s failed: %v, output: %s", ytDlpName, err, string(out))
		srv.logger.Error(errOut.Error())
		// Include stderr in error message
		return nil, errOut
	}

	outString := string(out)

	// Find JSON start
	start := strings.Index(outString, "{")
	if start == -1 {
		err := errors.New("no JSON found in yt-dlp output")
		srv.logger.Error(err.Error())
		return nil, err
	}

	// Extract JSON substring
	jsonData := outString[start:]

	var info = &dto.YouTubeInfo{}
	if err := json.NewDecoder(bytes.NewReader([]byte(jsonData))).Decode(info); err != nil {
		errOut := fmt.Errorf("no JSON found: %v", err)
		srv.logger.Error(errOut.Error())
		return nil, errOut
	}

	return srv.mappers.VideoInfoToDomain(info), nil
}
