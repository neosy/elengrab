CREATE TABLE IF NOT EXISTS files (
    -- Unique file identifier (UUID)
    file_id TEXT PRIMARY KEY,

    -- Associated user identifier (UUID)
    user_id TEXT NULL,

    -- Status
    file_status TEXT NOT NULL DEFAULT 'new', -- new, pending, working, done, failed

    -- Youtube URL
    youtube_url TEXT NOT NULL,

    -- Title from youtube
    youtube_title TEXT NOT NULL,

    -- Youtube channel ID
    youtube_channel_id TEXT NULL,

    -- Original file name
    file_name TEXT NOT NULL,
    
    -- File extension
    ext TEXT NOT NULL,
    
    -- Full file name (file_name + ext)
    full_name TEXT NOT NULL,

    -- File size (byte)
    file_size INTEGER NULL,

    -- Fast partial file hash (combined hash of multiple sampled blocks; not a full-file checksum)
    partial_hash TEXT NULL,
    
    -- Human-readable safe full name
    safe_readable_full_name TEXT NOT NULL,

    -- Media metadata as JSON (codecs, resolution, etc.)
    media_info TEXT NULL,

    -- Error message
    error_message TEXT NULL,
    
    -- Record creation timestamp, set automatically
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Record update timestamp, set automatically
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Record delete timestamp
    deleted_at DATETIME NULL,

    -- Foreign key linking files to user
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS files_created_at_idx
ON files(created_at);

CREATE INDEX IF NOT EXISTS files_partial_hash_idx
ON files(partial_hash);

CREATE INDEX files_deleted_at_null_idx
ON files(deleted_at)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS files_user_id_idx
ON files(user_id);
