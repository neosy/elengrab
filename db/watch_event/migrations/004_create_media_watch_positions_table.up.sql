CREATE TABLE media_watch_positions (
    -- Identifier of the watched media (UUID)
    download_id TEXT NOT NULL,

    -- Associated user identifier (UUID)
    user_id TEXT NOT NULL,

    -- User session identifier (UUID)
    session_id TEXT NOT NULL,

    -- Last saved playback position in milliseconds
    position_ms INTEGER NOT NULL,

    -- Record creation timestamp, set automatically
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Record last update timestamp, set automatically
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (download_id, user_id, session_id)
);
