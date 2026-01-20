-- Add column job_id - ID of the job currently processing the task
ALTER TABLE download_tasks ADD COLUMN job_id TEXT NULL;
