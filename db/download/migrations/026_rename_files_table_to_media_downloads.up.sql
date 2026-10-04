ALTER TABLE files RENAME TO media_downloads;

ALTER TABLE media_downloads RENAME COLUMN file_id TO download_id;

DROP INDEX IF EXISTS files_created_at_idx;
DROP INDEX IF EXISTS files_partial_hash_idx;
DROP INDEX IF EXISTS files_deleted_at_null_idx;
DROP INDEX IF EXISTS files_user_id_idx;
DROP INDEX IF EXISTS files_downloaded_created_sort_idx;
DROP INDEX IF EXISTS files_media_title_idx;

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