-- For new databases, define download_code as TEXT NOT NULL UNIQUE without a separate index.

-- Add download_code as NOT NULL with an empty default for existing records.
ALTER TABLE media_downloads ADD COLUMN download_code TEXT NOT NULL DEFAULT '';

UPDATE media_downloads
SET download_code = lower(hex(randomblob(8)));

-- Add a unique index because UNIQUE cannot be added to an existing column.
CREATE UNIQUE INDEX media_downloads_download_code_idx
ON media_downloads(download_code);