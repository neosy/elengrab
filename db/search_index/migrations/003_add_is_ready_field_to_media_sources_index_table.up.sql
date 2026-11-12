-- Add is_ready column to media_sources_index table
ALTER TABLE media_sources_index ADD COLUMN is_ready BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE media_sources_index SET is_ready = TRUE;