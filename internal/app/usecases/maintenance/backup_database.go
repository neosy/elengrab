package maintenance

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	backupname "github.com/neosy/elengrab/internal/infrastructure/sqlite/backup_name"
	nfile "github.com/neosy/elengrab/internal/pkg/filex"
)

func (m *maintenance) BackupDatabase(ctx context.Context) error {
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
		err := os.MkdirAll(backupDir, 0o755)
		if err != nil {
			m.logger.Error("Failed to make backup dir", "dir", backupDir, "error", err)
			return err
		}
	}

	for _, schema := range m.repositories.Schemas() {
		err = m.backupDB(schema.DBName(), backupDir)
		if err != nil {
			return err
		}
	}

	for _, schema := range m.repositories.Schemas() {
		err = m.cleanupOldBackups(schema.DBName(), backupDir)
		if err != nil {
			m.logger.Error("Failed cleanup old backups", "error", err)
			return err
		}
	}

	return nil
}

func (m *maintenance) backupDB(dbName string, backupsDir string) error {
	filePath := backupname.Path(dbName, backupsDir)

	if err := m.repositories.Backup(dbName, filePath); err != nil {
		m.logger.Error("Failed backup database", "error", err)
		return err
	}

	m.logger.Info(
		"Database backup completed",
		"dbName", dbName,
		"path", filePath,
	)

	return nil
}

func (m *maintenance) cleanupOldBackups(dbName string, backupsDir string) error {
	entries, err := os.ReadDir(backupsDir)
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
		if strings.HasPrefix(name, backupname.Prefix(dbName)) && strings.HasSuffix(name, backupname.FileExt) {
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
