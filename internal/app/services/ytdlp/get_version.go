package ytdlpsrv

import "context"

func (srv *YtDlpService) GetVersion(ctx context.Context) (string, error) {
	return srv.downloader.GetVersion(ctx)
}
