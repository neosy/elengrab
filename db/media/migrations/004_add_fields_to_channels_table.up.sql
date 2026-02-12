BEGIN;

-- Add columns for channel URL and title.
ALTER TABLE youtube_channels ADD COLUMN channel_url TEXT NOT NULL DEFAULT '';
ALTER TABLE youtube_channels ADD COLUMN channel_title TEXT NOT NULL DEFAULT '';

-- Create a new table with the additional fields.
CREATE TABLE IF NOT EXISTS youtube_channels_new (
    -- Unique ID for the channel
    channel_id TEXT PRIMARY KEY,

    -- Site URL
    channel_url TEXT NOT NULL,

    -- Title of the channel
    channel_title TEXT NOT NULL,

    -- URL of the image channel avatar
    image_url TEXT NOT NULL,

    -- Raw image data (binary)
    image_raw BLOB NOT NULL,

    -- Format of the image (jpg, png, webp)
    image_format TEXT NOT NULL,

    -- Record creation timestamp, set automatically
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Record update timestamp, set automatically
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Copy data from the old table to the new one.
INSERT INTO youtube_channels_new (channel_id, channel_url, channel_title, image_url, image_raw, image_format, created_at, updated_at)
SELECT channel_id, channel_url, channel_title, image_url, image_raw, image_format, created_at, updated_at FROM youtube_channels;

-- Drop the old table and rename the new one to the original name.
DROP TABLE youtube_channels;
ALTER TABLE youtube_channels_new RENAME TO youtube_channels;

COMMIT;
