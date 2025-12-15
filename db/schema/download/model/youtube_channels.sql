CREATE TABLE IF NOT EXISTS youtube_channels (
    -- Unique ID for the channel
    channel_id TEXT PRIMARY KEY,

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
)