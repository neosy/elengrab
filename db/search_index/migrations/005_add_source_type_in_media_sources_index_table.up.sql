-- Source of the media record: a remote URL or a file uploaded by the user
-- url, upload
ALTER TABLE media_sources_index ADD COLUMN source_type TEXT NOT NULL DEFAULT 'url';
