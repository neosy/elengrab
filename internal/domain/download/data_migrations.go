package ddownload

import (
	"time"
)

type DataMigration struct {
	// Unique identifier of the migration (e.g. "backfill_user_status")
	MigrationID string

	// Optional human-readable description of what this migration does
	Description *string

	// Timestamp when this migration record was created (i.e. when migration was applied)
	CreatedAt time.Time
}
