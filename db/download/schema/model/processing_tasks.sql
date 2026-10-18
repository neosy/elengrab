CREATE TABLE IF NOT EXISTS processing_tasks (
    -- Unique ID for the task
    task_id TEXT PRIMARY KEY,

    -- ID of the media record being processed
    download_id TEXT NOT NULL,

    -- Task status: new, pending, working, done, failed
    task_status TEXT NOT NULL DEFAULT 'new',

    -- Path to the staged uploaded file
    source_path TEXT NOT NULL,

    -- ID of the worker currently processing the task
    worker_id INT NULL,

    -- ID of the job currently processing the task
    job_id TEXT NULL,

    -- Task creation timestamp
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Last update timestamp
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Foreign key to media records
    FOREIGN KEY (download_id) REFERENCES media_downloads(download_id) ON DELETE CASCADE
);

-- Index for faster querying by download
CREATE INDEX IF NOT EXISTS processing_tasks_download_id_idx
ON processing_tasks(download_id);