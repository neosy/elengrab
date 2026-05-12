-- TODO: Rename to `media_downloads`. 
-- New table `media_sources` will be added as parent (one source → multiple files).
CREATE TABLE IF NOT EXISTS files (
    -- Unique file identifier (UUID)
    file_id TEXT PRIMARY KEY,

    -- Associated user identifier (UUID)
    user_id TEXT NULL,

    -- Status
    file_status TEXT NOT NULL DEFAULT 'new', -- new, pending, working, done, failed

    -- Media URL
    media_url TEXT NOT NULL,

    -- Title media
    media_title TEXT NOT NULL,

    -- MediaTitleLower in lowercase for efficient case-insensitive searches
    media_title_lower TEXT NOT NULL,

    -- Description media
    media_description TEXT NULL,

    -- Channel ID
    channel_id TEXT NULL,

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

    -- Downloaded timestamp
    downloaded_at DATETIME NULL,
    
    -- Record creation timestamp, set automatically
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Record update timestamp, set automatically
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Record delete timestamp
    deleted_at DATETIME NULL
);

-- Index for faster querying by creation date
CREATE INDEX IF NOT EXISTS files_created_at_idx
ON files(created_at);

-- Create index for partial_hash field.
CREATE INDEX IF NOT EXISTS files_partial_hash_idx
ON files(partial_hash);

-- Create index for deleted_at field where it is null.
-- This allows us to query only non-deleted records efficiently.
CREATE INDEX files_deleted_at_null_idx
ON files(deleted_at)
WHERE deleted_at IS NULL;

-- Create index for user_id field
CREATE INDEX IF NOT EXISTS files_user_id_idx
ON files(user_id);

-- Create index for sorting by download or update time, prioritizing downloads if available.
CREATE INDEX files_downloaded_created_sort_idx
ON files(COALESCE(downloaded_at, created_at) DESC);

-- Create index for youtube_title, media_title_lower fields
CREATE INDEX IF NOT EXISTS files_media_title_idx
ON files(youtube_title);
CREATE INDEX IF NOT EXISTS files_media_title_lower_idx
ON files(media_title_lower);
