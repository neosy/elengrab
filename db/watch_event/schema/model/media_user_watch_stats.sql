CREATE TABLE media_user_watch_stats (
    -- Identifier of the watched media (UUID)
    download_id   TEXT NOT NULL,

    -- Associated user identifier (UUID)
    -- Use '00000000-0000-0000-0000-000000000000' for anonymous users.
    user_id       TEXT NOT NULL,

    -- Number of completed views
    views         INTEGER NOT NULL DEFAULT 0,

    -- Record update timestamp, set automatically
    updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (download_id, user_id)
);