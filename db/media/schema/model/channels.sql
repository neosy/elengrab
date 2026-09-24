CREATE TABLE IF NOT EXISTS channels (
    -- Internal channel identifier
    channel_id TEXT PRIMARY KEY,

    -- External channel identifier
    external_id TEXT NOT NULL,

    -- Platform identifier
    platform TEXT NOT NULL,

    -- Site URL
    channel_url TEXT NOT NULL,

    -- Host from which the platform was detected
    host TEXT NOT NULL,
   
    -- Channel username
    username TEXT NOT NULL DEFAULT '',

    -- Channel URL based on the username
    username_url TEXT NOT NULL DEFAULT '',

   -- Title of the channel
    channel_title TEXT NOT NULL,

    -- URL of the image channel avatar
    image_url TEXT,

    -- Raw image data (binary)
    image_raw BLOB,

    -- Format of the image (jpg, png, webp)
    image_format TEXT,

    -- Record creation timestamp, set automatically
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Record update timestamp, set automatically
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Ensure external channel IDs are unique within each platform.
CREATE UNIQUE INDEX channels_platform_external_id_uidx
ON channels(platform, external_id);
