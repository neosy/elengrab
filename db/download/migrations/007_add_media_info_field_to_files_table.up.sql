-- Add column media_info - Media metadata as JSON (codecs, resolution, etc.)
ALTER TABLE files ADD COLUMN media_info TEXT NULL;
