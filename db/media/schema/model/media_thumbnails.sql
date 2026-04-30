-- Stores thumbnail metadata for media files
CREATE TABLE IF NOT EXISTS media_thumbnails (
    -- Unique identifier for the thumbnail record
    thumb_id TEXT PRIMARY KEY,

    -- Reference to the parent media entity
    media_id TEXT NOT NULL,

    -- Thumbnail variant type (e.g. small, medium, large, original)
    variant TEXT NOT NULL,

    -- Thumbnail generation version
    version INTEGER NOT NULL DEFAULT 1,

    -- Thumbnail width in pixels
    width INTEGER NULL,

    -- Thumbnail height in pixels
    height INTEGER NULL,

    -- Image format (e.g. jpg, png, webp)
    format TEXT NOT NULL,

    -- Source of the thumbnail (youtube, generated, upload, etc.)
    -- youtube | vimeo | external | video_frame | generated | upload
    source_type TEXT NOT NULL,

    -- Optional external identifier of the source
    source_id TEXT NULL,

    -- Optional external source URL used to derive the thumbnail
    source_url TEXT NULL,

    -- Stable object key used to resolve file location in storage backend (FS/S3/etc.)
    storage_key TEXT NOT NULL UNIQUE,

    -- Flag indicating whether this is the primary thumbnail (0/1)
    is_primary INTEGER NOT NULL DEFAULT 0
    CHECK (is_primary IN (0, 1)),

    -- Record creation timestamp
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Record last update timestamp
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Ensure only one variant per media
CREATE UNIQUE INDEX media_thumbnails_media_id_variant_version_uidx
ON media_thumbnails(media_id, variant, version);

-- Ensure only one primary thumbnail per media
CREATE UNIQUE INDEX media_thumbnails_media_primary_uidx
ON media_thumbnails(media_id)
WHERE is_primary = 1;

-- Ensures that each storage object key is unique across all thumbnail records,
-- preventing duplicate references to the same stored file
CREATE UNIQUE INDEX media_thumbnails_storage_key_uidx
ON media_thumbnails(storage_key);