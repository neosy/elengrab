-- Add columns
ALTER TABLE files ADD COLUMN media_title_lower TEXT NOT NULL DEFAULT '';

-- Create index for youtube_title, media_title_lower fields
CREATE INDEX IF NOT EXISTS files_media_title_idx
ON files(youtube_title);
CREATE INDEX IF NOT EXISTS files_media_title_lower_idx
ON files(media_title_lower);
