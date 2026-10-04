package ddownload

import (
	"time"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type DataMigration struct {
	// Unique identifier of the migration (e.g. "backfill_user_status")
	MigrationID string

	//Migration execution type (e.g. "required", "deferred")
	MigrationType dtypes.MigrationType

	// Optional human-readable description of what this migration does
	Description *string

	// Timestamp when this migration record was created (i.e. when migration was applied)
	CreatedAt time.Time
}
