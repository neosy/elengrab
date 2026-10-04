ALTER TABLE download_tasks RENAME COLUMN file_id TO download_id;
ALTER TABLE download_tasks RENAME COLUMN youtube_url TO media_url;

DROP INDEX IF EXISTS download_tasks_file_id_idx;

-- Create an index on the download_id column for faster querying by download.
CREATE INDEX IF NOT EXISTS download_tasks_download_id_idx
ON download_tasks(download_id);