CREATE TABLE media_sources_index (
    -- Identifier of the watched media (UUID)
    download_id TEXT NOT NULL PRIMARY KEY,

    -- User identifier
    user_id TEXT NULL,

    -- Title media
    title TEXT NOT NULL,

    -- Media title in lowercase for efficient case-insensitive searches
    title_lower TEXT NOT NULL,

    -- Description media
    description TEXT NULL,

    -- Description media in lowercase for efficient case-insensitive searches
    description_lower TEXT NOT NULL,

    -- Visibility access level for media (public, authenticated or private)
    visibility TEXT NOT NULL,

    -- Number of completed views
    views INTEGER NOT NULL DEFAULT 0,

    -- Media source creation timestamp
    source_created_at DATETIME NOT NULL,

    -- Record delete timestamp
    deleted_at DATETIME NULL
);

CREATE INDEX media_sources_index_source_created_at_idx
    ON media_sources_index (source_created_at);

CREATE INDEX media_sources_index_visibility_source_created_at_idx
    ON media_sources_index (visibility, source_created_at);

CREATE INDEX media_sources_index_visibility_views_idx
    ON media_sources_index (visibility, views);

CREATE INDEX media_sources_index_user_id_source_created_at_idx
    ON media_sources_index (user_id, source_created_at);