package downloader

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

const (
	// missingFileRetentionPeriod defines how long a missing download is retained
	// in soft-deleted state before it becomes eligible for permanent deletion.
	missingFileRetentionPeriod = 24 * time.Hour

	// moveUnmatchedFileRetentionPeriod defines how long unmatched downloads
	// should remain in the downloads directory before being moved.
	moveUnmatchedDownloadRetentionPeriod = 24 * time.Hour

	// lostDirName defines the name of the subdirectory where orphaned
	// or unmatched downloads from the downloads folder are moved.
	lostDirName = ".lost"
)

// DeleteMissingDownloads removes downloads from the database that no longer exist
// in the downloads directory, and moves any orphaned downloads found on disk
// into the "lost" folder for further inspection.
func (uc *Downloader) DeleteMissingDownloads(ctx context.Context, enableMoveUnmatchedDownloads bool) error {
	err := uc.deleteMissingDownloads(ctx)
	if err != nil {
		uc.logger.Error("Failed to delete missing downloads", "error", err)
		return err
	}

	if enableMoveUnmatchedDownloads {
		err = uc.moveUnmatchedDownloads(ctx)
		if err != nil {
			uc.logger.Error("Failed to move unmatched downloads", "error", err)
			return err
		}
	}

	return nil
}

// deleteMissingDownloads checks for downloads marked as Done,
// soft-deletes them if missing, and then either restores or permanently deletes
// records based on whether the download exists after the retention period.
func (uc *Downloader) deleteMissingDownloads(ctx context.Context) error {
	// Select all records with status Done
	downloads, err := uc.download.GetByStatus(ctx, dtypes.MediaDownloadStatusDone)
	if err != nil {
		return err
	}

	// Mark records for deletion if the download is missing
	for _, download := range downloads {
		exists, _ := uc.downloadsStorage.Exists(download.FileFullName)
		if !exists {
			err := uc.download.SoftDelete(ctx, download.DownloadID)
			if err == nil {
				uc.logger.Debug("Soft deleting file", "download_id", download.DownloadID, "fileName", download.FileFullName)
				uc.broadcastDownloadDelete(download.UserID, download.DownloadID)
			}
		}
	}

	// Select all records previously marked for deletion
	deletedDownloads, err := uc.download.GetDeleted(ctx, nil, nil)
	if err != nil {
		return err

	}

	// Restore records if the download was found
	// Permanently delete records if the download is still missing after the grace period
	for _, download := range deletedDownloads {
		if download.FileFullName == "" {
			continue
		}
		if exists, _ := uc.downloadsStorage.Exists(download.FileFullName); exists {
			err := uc.download.Restore(ctx, download.DownloadID)
			if err != nil {
				uc.logger.Warn("Failed to restore", "error", err)
				continue
			}
			uc.logger.Debug("Restoring media download in database", "downloadID", download.DownloadID, "fileName", download.FileFullName)
			continue
		}

		if time.Until(*download.DeletedAt) >= missingFileRetentionPeriod {
			err := uc.download.HardDelete(ctx, download.DownloadID)
			if err != nil {
				uc.logger.Warn("Failed to hard delete", "downloadID", download.DownloadID, "error", err)
				continue
			}
			uc.deleteThumbnails(ctx, download)
			uc.logger.Debug("Hard deleting media download from database", "downloadID", download.DownloadID, "fileName", download.FileFullName)
			continue
		}
	}

	return nil
}

// moveUnmatchedDownloads scans the downloads directory and moves all downloads
// that do not exist in the database into the "lost" subdirectory.
// This helps clean up orphaned downloads that were downloaded but not recorded.
func (uc *Downloader) moveUnmatchedDownloads(ctx context.Context) error {
	// Load all known filenames from DB (including deleted)
	names, err := uc.download.GetAllFullNames(ctx, true)
	if err != nil {
		return err
	}

	// Ensure "lost" directory exists
	lostDir := filepath.Join(uc.downloadsStorage.BasePath(), lostDirName)
	if err := os.MkdirAll(lostDir, 0755); err != nil {
		uc.logger.Warn("Failed to create lost dir", "error", err)
		return errorx.Errorf("failed to create lost dir: %w", err)
	}

	now := time.Now().UTC()

	// Read all entries from the downloads directory
	err = filepath.WalkDir(uc.downloadsStorage.MediaPath(),
		func(path string, file fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if file.IsDir() {
				return nil
			}

			fileName := file.Name()

			if strings.HasSuffix(fileName, ".part") {
				return nil
			}

			info, err := file.Info() // fs.FileInfo
			if err != nil {
				return nil
			}

			if now.Before(info.ModTime().UTC().Add(moveUnmatchedDownloadRetentionPeriod)) {
				return nil
			}

			// Move file if it's not present in DB
			if _, exists := names[fileName]; !exists {
				src := path
				dst := filepath.Join(uc.downloadsStorage.BasePath(), lostDirName, fileName)

				if err := os.Rename(src, dst); err != nil {
					uc.logger.Warn("Failed to move file", "file", fileName, "error", err)
					return fmt.Errorf("failed to move file %s: %w", fileName, err)
				}

				uc.logger.Debug("Moved file to lost directory", "file", fileName)
			}

			return nil
		},
	)
	if err != nil {
		uc.logger.Warn("Failed to read downloads dir", "error", err)
		return err
	}

	return nil
}
