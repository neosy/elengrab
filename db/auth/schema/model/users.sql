CREATE TABLE IF NOT EXISTS users (
    -- Unique user identifier (UUID)
    user_id TEXT PRIMARY KEY,

    -- Username/login, must be unique
    login TEXT NOT NULL UNIQUE,

    -- User email address, optional
    email TEXT NULL,

    -- User is guest: 1 = guest, 0 = not guest
    is_guest INTEGER NOT NULL DEFAULT 0,
    
    -- Active status: 1 = active, 0 = inactive
    is_active INTEGER NOT NULL DEFAULT 1,

    -- Record creation timestamp, set automatically
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Record update timestamp, set automatically
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for fast lookup by login (username)
CREATE INDEX IF NOT EXISTS users_login_idx
ON users(login);

-- Index for fast lookup by email
CREATE INDEX IF NOT EXISTS users_email_idx
ON users(email);
