CREATE TABLE media_watch_chunks (
    -- Identifier of the watched media (UUID)
    download_id   TEXT NOT NULL,

    -- Associated user identifier (UUID)
    -- Use '00000000-0000-0000-0000-000000000000' for anonymous users.
    user_id       TEXT NOT NULL,

    -- Zero-based index of the 1000ms media chunk
    chunk_index   INTEGER NOT NULL,

    -- How many times this chunk was watched
    qty           INTEGER NOT NULL DEFAULT 1,

    -- Record creation timestamp, set automatically
    created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (download_id, user_id, chunk_index)
);
