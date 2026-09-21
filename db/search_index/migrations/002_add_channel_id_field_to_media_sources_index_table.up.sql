ALTER TABLE media_sources_index ADD COLUMN channel_id TEXT NULL;

CREATE INDEX IF NOT EXISTS media_sources_index_channel_id_source_created_at_idx
    ON media_sources_index (channel_id, source_created_at);
