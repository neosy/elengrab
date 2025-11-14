-- Add columns
ALTER TABLE files ADD COLUMN file_status TEXT NOT NULL DEFAULT 'new';  -- new, pending, working, done, failed
ALTER TABLE files ADD COLUMN youtube_url TEXT NOT NULL DEFAULT '';
ALTER TABLE files ADD COLUMN error_message TEXT NULL;

-- Rename columns
ALTER TABLE files RENAME COLUMN title TO youtube_title;

-- Update all existing rows to have status 'done'
UPDATE files SET file_status = 'done';
