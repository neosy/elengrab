-- Add 'downloaded_at' column to the 'files' table
ALTER TABLE files ADD COLUMN downloaded_at DATETIME NULL;

-- Create index for sorting by download or update time, prioritizing downloads if available.
CREATE INDEX files_downloaded_created_sort_idx
ON files(COALESCE(downloaded_at, created_at) DESC);

-- Update the 'downloaded_at' column with the value of 'updated_at'
-- where it is null and file_status is 'done'
UPDATE files SET downloaded_at = created_at 
WHERE downloaded_at IS NULL AND file_status='done';