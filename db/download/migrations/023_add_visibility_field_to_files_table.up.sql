-- Add Visibility column - access level for media (public, authenticated or private)
ALTER TABLE files ADD COLUMN visibility TEXT NOT NULL DEFAULT 'private';

-- Set visibility to public for media (system media: user_id is zero UUID or NULL)
UPDATE files
SET visibility = 'public'
WHERE user_id = '00000000-0000-0000-0000-000000000000'
    OR user_id IS NULL;

-- Set visibility to authenticated for media
UPDATE files
SET visibility = 'authenticated'
WHERE user_id <> '00000000-0000-0000-0000-000000000000';
