-- Link Clicks Tracking Table
CREATE TABLE IF NOT EXISTS link_clicks
(
	-- Unique identifier for the click event
    link_click_id TEXT PRIMARY KEY,

    -- ID of the link that was clicked
    link_id TEXT NOT NULL,

    -- The IP address from which the link was accessed
	ip_address TEXT NOT NULL DEFAULT '',

    -- Full short URL, including domain, e.g., "https://s.nhub.ru/abc123"
    short_url TEXT NOT NULL,

	-- User ID who clicked the link (nullable, if not logged in or unknown)'
    clicked_by TEXT NULL,

    -- Timestamp of the click event
	clicked_at DATETIME NOT NULL,

	-- User agent or device info (optional for tracking purposes)
    user_agent TEXT NULL,

	-- Referrer URL (optional, if available)
    referrer TEXT NULL,

    -- Timestamp when the event was created
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Foreign key linking to the 'links' table, ensuring referential integrity
    FOREIGN KEY (link_id) REFERENCES links(link_id) ON DELETE CASCADE
);

-- Creating index for link_id column
CREATE INDEX IF NOT EXISTS link_clicks_link_id_idx
ON link_clicks(link_id);

-- Creating index for created_at column
CREATE INDEX IF NOT EXISTS link_clicks_created_at_idx
ON link_clicks(created_at DESC);
