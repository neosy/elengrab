-- Channel username
ALTER TABLE channels ADD COLUMN username TEXT NOT NULL DEFAULT '';
-- Channel URL based on the username
ALTER TABLE channels ADD COLUMN username_url TEXT NOT NULL DEFAULT '';
