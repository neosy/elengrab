package dlstate

import "context"

func (uc *DownloadStateCache) Transaction(fn func(context.Context) error) error {
	return uc.stateCacheRep.Transaction(fn)
}
