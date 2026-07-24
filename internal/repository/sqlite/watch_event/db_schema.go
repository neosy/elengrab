package watchevent

import (
	"github.com/neosy/elengrab/internal/ports/persistence"
	sqlitetypes "github.com/neosy/elengrab/internal/repository/sqlite/types"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/table_names"
)

func NewDBSchema(name string) persistence.DBSchema {
	return sqlitetypes.NewDBSchema(name, tablenames.TableNames())
}
