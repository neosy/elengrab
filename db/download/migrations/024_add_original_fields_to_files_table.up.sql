-- Add original title from the media source
ALTER TABLE files ADD COLUMN media_title_original TEXT;

-- Add original description from the media source
ALTER TABLE files ADD COLUMN media_description_original TEXT;

UPDATE files SET
    media_title_original = media_title,
    media_description_original = media_description;