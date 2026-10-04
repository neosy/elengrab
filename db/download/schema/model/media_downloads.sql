-- table `media_downloads` will be added as parent (one source → multiple files).
CREATE TABLE IF NOT EXISTS media_downloads (
    -- Unique media download identifier (UUID)
    download_id TEXT PRIMARY KEY,

    -- Unique short download code
    download_code TEXT NOT NULL UNIQUE,

    -- Associated user identifier (UUID)
    user_id TEXT NULL,

    -- Status
    file_status TEXT NOT NULL DEFAULT 'new', -- new, pending, working, done, failed

    -- Media URL
    media_url TEXT NOT NULL,

    -- Original title from the media source
    media_title_original TEXT NOT NULL,

    -- Title media
    media_title TEXT NOT NULL,

    -- Original description from the media source
    media_description_original TEXT NULL,

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

    -- Visibility access level for media (public, authenticated or private)
    visibility TEXT NOT NULL DEFAULT 'private',

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
CREATE INDEX IF NOT EXISTS media_downloads_created_at_idx
ON media_downloads(created_at);

-- Create index for partial_hash field.
CREATE INDEX IF NOT EXISTS media_downloads_partial_hash_idx
ON media_downloads(partial_hash);

-- Create index for deleted_at field where it is null.
-- This allows us to query only non-deleted records efficiently.
CREATE INDEX IF NOT EXISTS media_downloads_deleted_at_null_idx
ON media_downloads(deleted_at)
WHERE deleted_at IS NULL;

-- Create index for user_id field
CREATE INDEX IF NOT EXISTS media_downloads_user_id_idx
ON media_downloads(user_id);

-- Create index for sorting by download or update time, prioritizing downloads if available.
CREATE INDEX IF NOT EXISTS media_downloads_downloaded_created_sort_idx
ON media_downloads(COALESCE(downloaded_at, created_at) DESC);

-- Create index for media_title fields
CREATE INDEX IF NOT EXISTS media_downloads_media_title_idx
ON media_downloads(media_title);