CREATE TABLE media_watch_events (
    -- Unique event identifier (UUID)
    event_id TEXT PRIMARY KEY,

    -- Identifier of the watched media (UUID)
    download_id TEXT NOT NULL,

    -- Associated user identifier (UUID)
    user_id TEXT NULL,

    -- Unique viewing session identifier (UUID)
    session_id TEXT NULL,

    -- Playback position in milliseconds
    position_ms INTEGER NOT NULL,

    -- Playback duration since the previous event in milliseconds
    interval_ms INTEGER NOT NULL,

    -- Record creation timestamp, set automatically
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX media_watch_events_download_id_created_idx
    ON media_watch_events (download_id, created_at);

CREATE INDEX media_watch_events_user_id_created_idx
    ON media_watch_events (user_id, created_at);