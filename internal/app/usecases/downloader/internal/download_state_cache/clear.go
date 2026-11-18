package dlstate

import "context"

func (uc *DownloadStateCache) Clear(ctx context.Context) error {
	return uc.stateCacheRep.Clear(ctx)
}
