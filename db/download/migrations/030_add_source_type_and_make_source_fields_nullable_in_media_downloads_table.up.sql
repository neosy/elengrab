-- Source of the media record: a remote URL or a file uploaded by the user
ALTER TABLE media_downloads ADD COLUMN source_type TEXT NOT NULL DEFAULT 'url'; -- url, upload

-- Uploaded files have no source URL
ALTER TABLE media_downloads ALTER COLUMN media_url DROP NOT NULL;