CREATE TABLE data_migrations (
    -- Unique identifier of the migration (e.g. "2026-05-03_backfill_user_status")
    migration_id TEXT PRIMARY KEY,

    -- Optional human-readable description of what this migration does
    description TEXT NULL,

    -- Timestamp when this migration record was created (i.e. when migration was applied)
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);