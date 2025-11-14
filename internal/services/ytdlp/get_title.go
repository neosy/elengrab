package ytdlpsrv

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// GetTitle
func (srv *YtDlpService) GetTitle(url string) (string, error) {
	cmd := exec.Command(srv.cmdPath, "--no-playlist", "-e", url)

	// Buffers to capture stdout and stderr
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errOut := fmt.Errorf("%s failed: %v, stderr: %s", ytDlpName, err, stderr.String())
		srv.logger.Error(errOut.Error())
		return "", errOut
	}

	title := strings.TrimSpace(out.String())
	if title == "" {
		err := fmt.Errorf("title not found")
		srv.logger.Error(err.Error())
		return "", err
	}

	return title, nil
}
