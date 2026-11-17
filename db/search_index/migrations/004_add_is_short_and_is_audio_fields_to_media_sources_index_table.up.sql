-- Whether the media is a short (short-form vertical video)
ALTER TABLE media_sources_index ADD COLUMN is_short BOOLEAN NOT NULL DEFAULT FALSE;

-- Whether the media is audio only (no video track)
ALTER TABLE media_sources_index ADD COLUMN is_audio BOOLEAN NOT NULL DEFAULT FALSE;