-- 1. Temporarily disable foreign key checks
PRAGMA foreign_keys = OFF;

-- 2. Create a new users table with UNIQUE constraint on email
CREATE TABLE IF NOT EXISTS users_new (
    -- Unique user identifier (UUID)
    user_id TEXT PRIMARY KEY,

    -- Username/login, must be unique
    login TEXT NOT NULL UNIQUE,

    -- User email address, optional
    email TEXT NULL UNIQUE,

    -- Active status: 1 = active, 0 = inactive
    is_active INTEGER NOT NULL DEFAULT 1,

    -- Record creation timestamp, set automatically
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Record update timestamp, set automatically
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 3. Copy existing data into the new table
INSERT INTO users_new (user_id, login, email, is_active, created_at, updated_at)
SELECT user_id, login, email, is_active, created_at, updated_at
FROM users;

-- 4. Create a temporary table and copy existing user_roles
CREATE TABLE user_roles_temp AS
SELECT *
FROM user_roles;

-- 5. Create a temporary table and copy existing user_sessions
CREATE TABLE user_sessions_temp AS
SELECT *
FROM user_sessions;

-- 6. Drop the old users table
DROP TABLE users;

-- 7. Rename the new table to users
ALTER TABLE users_new RENAME TO users;

-- 8. Recreate indexes
-- Index for fast lookup by login (username)
CREATE INDEX IF NOT EXISTS users_login_idx
ON users(login);

-- Index for fast lookup by email
CREATE INDEX IF NOT EXISTS users_email_idx
ON users(email);

-- 9. Restore data from the temporary table
DELETE FROM user_sessions;
INSERT INTO user_sessions SELECT * FROM user_sessions_temp;
DELETE FROM user_roles;
INSERT INTO user_roles SELECT * FROM user_roles_temp;

-- 10. Drop the temporary tables
DROP TABLE user_roles_temp;
DROP TABLE user_sessions_temp;

-- 11. Re-enable foreign key checks
PRAGMA foreign_keys = ON;