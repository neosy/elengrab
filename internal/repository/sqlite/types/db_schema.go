package types

import (
	"fmt"
	"path/filepath"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

// Database represents the name of a database.
type dbSchema struct {
	name       string
	tableNames []string
}

func NewDBSchema(name string, tableNames []string) persistence.DBSchema {
	return &dbSchema{
		name:       name,
		tableNames: tableNames,
	}
}

// Name returns the string representation of a database name.
func (s *dbSchema) DBName() string {
	return s.name
}

// Path returns the database path inside the given directory
func (s *dbSchema) Path(dir string) string {
	return filepath.Join(dir, s.FileName())
}

func (s *dbSchema) FileName() string {
	return fmt.Sprintf("%s.db", s.name)
}

func (s *dbSchema) TableNames() []string {
	return s.tableNames
}
