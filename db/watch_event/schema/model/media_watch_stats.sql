CREATE TABLE media_watch_stats (
    -- Identifier of the watched media (UUID)
    download_id   TEXT NOT NULL PRIMARY KEY,

    -- Number of completed views
    views         INTEGER NOT NULL DEFAULT 0,

    -- Record update timestamp, set automatically
    updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);