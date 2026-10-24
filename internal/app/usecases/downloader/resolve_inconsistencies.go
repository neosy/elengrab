package downloader

import (
	"context"
	"fmt"
	"time"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
)

// ResolveInconsistencies removes index entries that have no matching download.
//
// Only entries older than orphanGracePeriod are checked, so entries that are
// still being created are not touched. A zero orphanGracePeriod disables
// this protection and checks ALL entries: use it only when no downloads
// are being created concurrently.
func (uc *downloader) ResolveInconsistencies(ctx context.Context, orphanGracePeriod time.Duration) error {
	if orphanGracePeriod < 0 {
		return fmt.Errorf("orphan grace period must not be negative, got %v", orphanGracePeriod)
	}

	if orphanGracePeriod == 0 {
		uc.logger.Warn("Resolving inconsistencies without grace period: all index entries will be checked")
	}

	var options *dtypes.QueryMediaOptions

	if orphanGracePeriod > 0 {
		options = new(dtypes.NewQueryMediaOptions())

		options.Filters.Append(
			dtypes.NewQueryFilter(
				dtypes.QueryFilterNameCreatedAt,
				dbutils.NewFilterCondition(time.Now().UTC().Add(-orphanGracePeriod), dbutils.FilterOperatorLt),
			))
	}

	ids, err := uc.searchIndex.GetDownloadIDsFromMediaSourceIndex(ctx, options)
	if err != nil {
		return fmt.Errorf("get download ids from index: %w", err)
	}

	deleted := 0

	for _, id := range ids {
		// Stop early if the caller cancelled the context.
		if err := ctx.Err(); err != nil {
			return err
		}

		exists, err := uc.download.ExistsByDownloadIDNoCache(ctx, id)
		if err != nil {
			return fmt.Errorf("check download %v exists: %w", id, err)
		}

		if exists {
			continue
		}

		err = uc.searchIndex.HardDeleteMediaDownload(ctx, id)
		if err != nil {
			return fmt.Errorf("hard delete media download %v: %w", id, err)
		}

		// Log every deletion: the operation is irreversible.
		uc.logger.Info("Deleted orphaned index entry", "download_id", id)

		deleted++
	}

	if deleted > 0 {
		uc.logger.Info("Inconsistencies resolved", "checked", len(ids), "deleted", deleted)
	}

	return nil
}
