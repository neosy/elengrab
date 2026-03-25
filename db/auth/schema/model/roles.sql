-- Roles
CREATE TABLE roles (
    -- Unique role identifier (can be a readable key like 'admin', 'guest')
    role_id TEXT PRIMARY KEY,

    -- Human-readable role name, must be unique across the system
    name TEXT NOT NULL UNIQUE,

    -- Optional description of role
    description TEXT NULL,

    -- Record creation timestamp, set automatically
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Record update timestamp, set automatically
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);