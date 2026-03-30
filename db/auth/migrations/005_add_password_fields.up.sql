-- User password hash, optional
ALTER TABLE users ADD COLUMN password_hash TEXT NULL;
-- Timestamp when the user's password was last updated
ALTER TABLE users ADD COLUMN password_updated_at DATETIME NULL;
-- Record delete timestamp
ALTER TABLE users ADD COLUMN deleted_at DATETIME NULL;