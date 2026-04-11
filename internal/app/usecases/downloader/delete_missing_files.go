package downloader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	nfile "github.com/neosy/elengrab/internal/pkg/file"
)

const (
	// missingFileRetentionPeriod defines how long a missing file is retained
	// in soft-deleted state before it becomes eligible for permanent deletion.
	missingFileRetentionPeriod = 24 * time.Hour

	// moveUnmatchedFileRetentionPeriod defines how long unmatched files
	// should remain in the downloads directory before being moved.
	moveUnmatchedFileRetentionPeriod = 24 * time.Hour

	// lostDirName defines the name of the subdirectory where orphaned
	// or unmatched files from the downloads folder are moved.
	lostDirName = ".lost"
)

// DeleteMissingFiles removes files from the database that no longer exist
// in the downloads directory, and moves any orphaned files found on disk
// into the "lost" folder for further inspection.
func (uc *Downloader) DeleteMissingFiles(ctx context.Context, enableMoveUnmatchedFiles bool) error {
	err := uc.deleteMissingFiles(ctx)
	if err != nil {
		uc.logger.Error("Failed to delete missing files", "error", err)
		return err
	}

	if enableMoveUnmatchedFiles {
		err = uc.moveUnmatchedFiles(ctx)
		if err != nil {
			uc.logger.Error("Failed to move unmatched files", "error", err)
			return err
		}
	}

	return nil
}

// deleteMissingFiles checks for files marked as Done,
// soft-deletes them if missing, and then either restores or permanently deletes
// records based on whether the file exists after the retention period.
func (uc *Downloader) deleteMissingFiles(ctx context.Context) error {
	// Select all records with status Done
	files, err := uc.file.GetByStatus(ctx, dtypes.FileStatusDone)
	if err != nil {
		return err
	}

	// Mark records for deletion if the file is missing
	for _, file := range files {
		filePath := filepath.Join(uc.downloadsDir, file.FullName)
		exists, _ := nfile.FileExists(filePath)
		if !exists {
			err := uc.file.SoftDelete(ctx, file.FileID)
			if err == nil {
				uc.logger.Debug("Soft deleting file", "file_id", file.FileID, "fileName", file.FullName)
				uc.broadcastFileDelete(file.UserID, file.FileID)
			}
		}
	}

	// Select all records previously marked for deletion
	deletedFiles, err := uc.file.GetDeleted(ctx, nil, nil)
	if err != nil {
		return err

	}

	// Restore records if the file was found
	// Permanently delete records if the file is still missing after the grace period
	for _, file := range deletedFiles {
		if file.FullName == "" {
			continue
		}
		filePath := filepath.Join(uc.downloadsDir, file.FullName)
		if exists, _ := nfile.FileExists(filePath); exists {
			err := uc.file.Restore(ctx, file.FileID)
			if err != nil {
				uc.logger.Warn("Failed to restore", "error", err)
				continue
			}
			uc.logger.Debug("Restoring file in database", "fileId", file.FileID, "fileName", file.FullName)
			continue
		}

		if time.Until(*file.DeletedAt) >= missingFileRetentionPeriod {
			err := uc.file.HardDelete(ctx, file.FileID)
			if err != nil {
				uc.logger.Warn("Failed to hard delete", "fileID", file.FileID, "error", err)
				continue
			}
			uc.logger.Debug("Hard deleting file from database", "fileId", file.FileID, "fileName", file.FullName)
			continue
		}
	}

	return nil
}

// moveUnmatchedFiles scans the downloads directory and moves all files
// that do not exist in the database into the "lost" subdirectory.
// This helps clean up orphaned files that were downloaded but not recorded.
func (uc *Downloader) moveUnmatchedFiles(ctx context.Context) error {
	// Read all entries from the downloads directory
	files, err := os.ReadDir(uc.downloadsDir)
	if err != nil {
		uc.logger.Warn("Failed to read downloads dir", "error", err)
		return err
	}

	// Load all known filenames from DB (including deleted)
	names, err := uc.file.GetAllFullNames(ctx, true)
	if err != nil {
		return err
	}

	// Convert list to a set for fast lookup
	var namesMap = make(map[string]struct{}, len(names))
	for _, name := range names {
		namesMap[name] = struct{}{}
	}

	// Ensure "lost" directory exists
	lostDir := filepath.Join(uc.downloadsDir, lostDirName)
	if err := os.MkdirAll(lostDir, 0755); err != nil {
		uc.logger.Warn("Failed to create lost dir", "error", err)
		return errorx.Errorf("failed to create lost dir: %w", err)
	}

	// Iterate over all files in downloads dir
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) == ".part" {
			continue
		}

		info, err := file.Info() // fs.FileInfo
		if err != nil {
			continue
		}

		if time.Now().UTC().Before(info.ModTime().UTC().Add(moveUnmatchedFileRetentionPeriod)) {
			continue
		}

		// Move file if it's not present in DB
		if _, exists := namesMap[file.Name()]; !exists {
			src := filepath.Join(uc.downloadsDir, file.Name())
			dst := filepath.Join(uc.downloadsDir, lostDirName, file.Name())

			if err := os.Rename(src, dst); err != nil {
				uc.logger.Warn("Failed to move file", "file", file.Name(), "error", err)
				return fmt.Errorf("failed to move file %s: %w", file.Name(), err)
			}

			uc.logger.Debug("Moved file to lost directory", "file", file.Name())
		}
	}

	return nil
}
