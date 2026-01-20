-- Add column deleted_at
ALTER TABLE files ADD COLUMN deleted_at DATETIME NULL;

-- Add index files_deleted_at_null_idx
CREATE INDEX files_deleted_at_null_idx
ON files(deleted_at)
WHERE deleted_at IS NULL;
