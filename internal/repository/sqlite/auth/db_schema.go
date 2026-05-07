package auth

import (
	"github.com/neosy/elengrab/internal/ports/persistence"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/auth/table_names"
	sqlitetypes "github.com/neosy/elengrab/internal/repository/sqlite/types"
)

func NewDBSchema(name string) persistence.DBSchema {
	return sqlitetypes.NewDBSchema(name, tablenames.TableNames())
}
