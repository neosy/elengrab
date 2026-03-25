-- User sessions
CREATE TABLE IF NOT EXISTS user_sessions (
    -- Unique session identifier (UUID)
    session_id TEXT PRIMARY KEY,

    -- Associated user identifier (UUID)
    user_id TEXT NOT NULL,

    -- Random session token stored in cookie
    session_token TEXT NOT NULL UNIQUE,

    -- Record creation timestamp, set automatically
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Session expiration timestamp
    expires_at DATETIME NOT NULL,

    -- Foreign key linking session to user
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Index for fast lookup of sessions by user
CREATE INDEX IF NOT EXISTS user_sessions_user_id_idx
ON user_sessions(user_id);

-- Index for fast lookup of sessions by token (used on login/auth)
CREATE INDEX IF NOT EXISTS user_sessions_token_idx
ON user_sessions(session_token);

-- Index for efficiently cleaning up expired sessions
CREATE INDEX IF NOT EXISTS user_sessions_expires_at_idx
ON user_sessions(expires_at);
