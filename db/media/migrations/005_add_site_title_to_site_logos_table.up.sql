BEGIN;

-- Add column for site title.
ALTER TABLE site_logos ADD COLUMN 
    site_title TEXT NOT NULL DEFAULT '';

-- Create a new table with the additional fields.
CREATE TABLE IF NOT EXISTS site_logos_new (
    -- Unique ID for the logo 
    logo_id TEXT PRIMARY KEY,

    -- Site URL
    site_url TEXT NOT NULL,

    -- Title of the site
    site_title TEXT NOT NULL,

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

-- Copy data from the old table to the new table.
INSERT INTO site_logos_new (logo_id, site_url, site_title, image_url, image_raw, image_format, created_at, updated_at)
SELECT logo_id, site_url, site_title, image_url, image_raw, image_format, created_at, updated_at FROM site_logos;

-- Drop the old table.
DROP TABLE site_logos;
-- Rename the new table to the old table name.
ALTER TABLE site_logos_new RENAME TO site_logos;

-- Ensure only one logo per site URL (create index after table rename).
CREATE UNIQUE INDEX site_logos_site_url_uidx
ON site_logos(site_url);

COMMIT;