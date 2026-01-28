PRAGMA foreign_keys = OFF;

ALTER TABLE download_tasks RENAME TO download_tasks_old;

/* delete broken links */
DELETE FROM download_tasks_old
WHERE file_id NOT IN (SELECT file_id FROM files);

CREATE TABLE download_tasks (
    -- Unique ID for the task
    task_id TEXT PRIMARY KEY,
    
    -- ID of the file to download
    file_id TEXT NOT NULL,
    
    -- Task status: pending, working, done, failed
    task_status TEXT NOT NULL DEFAULT 'new',

    -- Youtube URL
    youtube_url TEXT NOT NULL,

	-- Youtube download options
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
    FOREIGN KEY (file_id) REFERENCES files(file_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS download_tasks_file_id_idx
ON download_tasks(file_id);

INSERT INTO download_tasks (
    task_id,
    file_id,
    task_status,
    youtube_url,
    options,
    worker_id,
    job_id,
    created_at,
    updated_at
)
SELECT
    task_id,
    file_id,
    task_status,
    youtube_url,
    options,
    worker_id,
    job_id,
    created_at,
    updated_at
FROM download_tasks_old;

DROP TABLE download_tasks_old;

PRAGMA foreign_keys = ON;
