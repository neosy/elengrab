package ytdlpsrv

import (
	"context"

	"github.com/neosy/elengrab/internal/app/services/internal/utils"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/executor"
)

func (srv *YtDlpService) GetVersion(ctx context.Context) (string, error) {
	return srv.downloader.GetVersion(ctx)
}

// GetVersion retrieves the version of the YtDlp executable.
func GetVersion(ctx context.Context) (string, error) {
	cmdPath, err := utils.ResolveCmdPath(consts.YtDlpName, "")
	if err != nil {
		return "", err
	}

	return executor.GetVersion(ctx, cmdPath)
}
