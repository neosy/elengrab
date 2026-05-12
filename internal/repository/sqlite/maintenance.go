package sqliterep

import "context"

func (r *Repositories) StartupMaintenance(ctx context.Context) error {
	return r.MediaDownload.FillEmptyMediaTitleLower(ctx)
}
