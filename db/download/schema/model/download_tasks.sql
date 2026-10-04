CREATE TABLE IF NOT EXISTS download_tasks (
    -- Unique ID for the task
    task_id TEXT PRIMARY KEY,
    
    -- ID of the file to download
    download_id TEXT NOT NULL,
    
    -- Task status: pending, working, done, failed
    task_status TEXT NOT NULL DEFAULT 'new',

    -- Media URL
    media_url TEXT NOT NULL,

	-- Media download options
	options TEXT NULL,

    -- ID of the worker currently processing the task
    worker_id INT NULL,

    -- ID of the job currently processing the task
    job_id TEXT NULL,
    
    -- Task creation timestamp
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Last update timestamp
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Foreign key to files
    FOREIGN KEY (download_id) REFERENCES media_downloads(download_id) ON DELETE CASCADE
);

-- Create an index on the download_id column for faster querying by download.
CREATE INDEX IF NOT EXISTS download_tasks_download_id_idx
ON download_tasks(download_id);