package ytdlpsrv

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/dto"
	dyoutubeinfo "github.com/neosy/elengrab/internal/domain/youtube_info"
)

// GetFormats retrieves and parses video formats for the given URL.
func (srv *YtDlpService) GetFormats(ctx context.Context, url string) (*dyoutubeinfo.YouTubeInfo, error) {
	// Prepare command to fetch video info in JSON
	cmd := exec.CommandContext(ctx, srv.cmdPath, "--no-playlist", "--no-warnings", "-J", url)

	// Capture combined stdout and stderr
	out, err := cmd.CombinedOutput()
	if err != nil {
		// The context was canceled
		if ctx.Err() != nil {
			return nil, fmt.Errorf("process canceled: %w", ctx.Err())
		}

		errOut := fmt.Errorf("%s failed: %v, output: %s", ytDlpName, err, string(out))
		// Include stderr in error message
		return nil, errOut
	}

	outStr := string(out)

	// Find JSON start
	start := strings.Index(outStr, "{")
	if start == -1 {
		err := errors.New("no JSON found in yt-dlp output")
		return nil, err
	}

	// Extract JSON substring
	jsonData := outStr[start:]

	var info = &dto.YouTubeInfo{}
	if err := json.NewDecoder(bytes.NewReader([]byte(jsonData))).Decode(info); err != nil {
		errOut := fmt.Errorf("no JSON found: %v", err)
		return nil, errOut
	}

	return srv.mappers.VideoInfoToDomain(info), nil
}
