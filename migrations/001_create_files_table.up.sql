-- golang-migrate/migrate
-- migrate -path ./migrations -database sqlite3://files.db up

CREATE TABLE IF NOT EXISTS files (
    -- Unique file identifier (UUID)
    file_id TEXT PRIMARY KEY,

    -- Title from youtube
    title TEXT NOT NULL,
    
    -- Original file name
    file_name TEXT NOT NULL,
    
    -- File extension
    ext TEXT NOT NULL,
    
    -- Full file name (file_name + ext)
    full_name TEXT NOT NULL,
    
    -- Human-readable safe full name
    safe_readable_full_name TEXT NOT NULL,
    
    -- Record creation timestamp, set automatically
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Record update timestamp, set automatically
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
