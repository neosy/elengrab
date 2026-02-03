CREATE TABLE IF NOT EXISTS site_logos (
    -- Unique ID for the logo 
    logo_id TEXT PRIMARY KEY,

    -- Site URL
    site_url TEXT NOT NULL,

    -- URL of the logo image
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

-- Ensure only one logo per site URL
CREATE UNIQUE INDEX site_logos_site_url_uidx
ON site_logos(site_url);