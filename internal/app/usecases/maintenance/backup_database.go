package maintenance

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/neosy/elengrab/pkg/nfile"
)

func (m *Maintenance) BackupDatabase(ctx context.Context) error {
	if m.databaseBackupsKeep <= 0 {
		return nil
	}

	backupDir := m.databaseBackupsDir
	exists, err := nfile.DirExists(backupDir)
	if err != nil {
		m.logger.Error("Failed to check backup dir", "dir", backupDir, "error", err)
		return err

	}
	if !exists {
		err := os.MkdirAll(backupDir, os.FileMode(0755))
		if err != nil {
			m.logger.Error("Failed to make backup dir", "dir", backupDir, "error", err)
			return err
		}
	}

	filename := fmt.Sprintf(
		"%s_%s.%s",
		m.prefixBackup(),
		time.Now().Format(backupDateFormat),
		backupFileExt,
	)

	path := filepath.Join(backupDir, filename)

	if err := m.database.Backup(path); err != nil {
		m.logger.Error("Failed backup database", "error", err)
		return err
	}

	m.logger.Info("Database backup completed", "path", path)

	if err := m.cleanupOldBackups(); err != nil {
		m.logger.Error("Failed cleanup old backups", "error", err)
		return err
	}

	return nil
}

func (m *Maintenance) prefixBackup() string {
	prefix := nfile.SanitizeFileName(m.appName)
	prefix = strings.ToLower(prefix)
	return prefix
}

func (m *Maintenance) cleanupOldBackups() error {
	entries, err := os.ReadDir(m.databaseBackupsDir)
	if err != nil {
		return err
	}

	files := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()

		// Filter by prefix and extension
		if strings.HasPrefix(name, m.prefixBackup()) && strings.HasSuffix(name, backupFileExt) {
			files = append(files, name)
		}
	}

	// Sort by name in reverse order (new files first)
	sort.Slice(files, func(i, j int) bool {
		return files[i] > files[j]
	})

	// Remove old files, leave keep latest
	for i, f := range files {
		if i+1 > m.databaseBackupsKeep {
			path := filepath.Join(m.databaseBackupsDir, f)

			if err := os.Remove(path); err != nil {
				m.logger.Warn("Failed to remove backup file", "file", path, "error", err)
				continue
			}

			m.logger.Debug("Remove backup file", "file", path)
		}
	}

	return nil
}
