-- Add column user_id
ALTER TABLE files ADD COLUMN user_id TEXT NULL;

CREATE INDEX IF NOT EXISTS files_user_id_idx
ON files(user_id);
