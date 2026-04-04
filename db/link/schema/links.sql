-- Links
CREATE TABLE IF NOT EXISTS links
(
    -- Unique identifier of the link
    link_id TEXT PRIMARY KEY,

    -- Original (long) URL to be shortened
    original_url TEXT NOT NULL,

    -- Generated short code used in the shortened URL, e.g., "abc123
    short_code TEXT NOT NULL,

    -- Full short URL, including domain, e.g., "https://s.nhub.ru/abc123"
    short_url TEXT NOT NULL,

    -- Indicates if the full short URL should be used for exact match
	is_match_short_url INTEGER NOT NULL DEFAULT 0,

    -- Maximum number of allowed clicks; nil means unlimited
    max_clicks INTEGER NULL,

    -- JSON array of user IDs allowed to access the link; nil means no restrictions
    allowed_user_ids TEXT NULL,

    -- JSON array of IP addresses allowed to access the link; nil means no restrictions
    allowed_ips TEXT NULL,

    -- Expiration date and time for the link; nil means no expiration
    expires_at DATETIME NULL,

    -- Timestamp when the link was created
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Timestamp when the link was last updated
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Timestamp when the link was soft-deleted
    deleted_at DATETIME NULL
);

-- Creating index for short_code column
CREATE INDEX IF NOT EXISTS links_short_code_idx
ON links(short_code);

-- Creating index for expires_at column
CREATE INDEX IF NOT EXISTS links_expires_at_idx
ON links(expires_at);

-- Creating index for created_at column
CREATE INDEX IF NOT EXISTS links_created_at_idx
ON links(created_at DESC);

-- Creating index for deleted_at column
CREATE INDEX IF NOT EXISTS links_deleted_at_idx
ON links(deleted_at);
