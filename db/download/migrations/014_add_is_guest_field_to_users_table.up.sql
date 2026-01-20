-- Add column is_guest
ALTER TABLE users ADD COLUMN is_guest INTEGER NOT NULL DEFAULT 0;

-- Update rocords
UPDATE users SET is_guest = 1;
