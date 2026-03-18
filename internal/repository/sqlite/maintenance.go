package sqliterep

import "context"

func (r *Repositories) StartupMaintenance(ctx context.Context) error {
	return r.File.FillEmptyMediaTitleLower(ctx)
}
