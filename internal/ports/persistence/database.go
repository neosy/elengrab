package persistence

import (
	"context"
	"fmt"
	"path/filepath"
)

// DBName represents the name of a database.
type DBName string

const (
	DBMainName  DBName = "elengrab"
	DBAuthName  DBName = "auth"
	DBMediaName DBName = "media"
)

// String returns the string representation of a DBName.
func (n DBName) String() string {
	return string(n)
}

// Path returns the database path inside the given directory
func (n DBName) Path(dir string) string {
	return filepath.Join(dir, n.FileName())
}

func (n DBName) FileName() string {
	return fmt.Sprintf("%s.db", n)
}

// Database is an interface that defines methods for interacting with databases.
type Database interface {
	Backup(dbName DBName, path string) error
	FlushWAL() error
	GetDBNames() []DBName
	StartupMaintenance(ctx context.Context) error
}
