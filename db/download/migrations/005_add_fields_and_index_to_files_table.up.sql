-- Add column partial_hash
ALTER TABLE files ADD COLUMN partial_hash TEXT NULL;

-- Add index files_created_at_idx
CREATE INDEX IF NOT EXISTS files_created_at_idx
ON files(created_at);

-- Add index files_partial_hash_idx
CREATE INDEX IF NOT EXISTS files_partial_hash_idx
ON files(partial_hash);