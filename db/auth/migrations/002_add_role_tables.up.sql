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

-- User roles
CREATE TABLE user_roles (
    -- Reference to the user (many-to-many relationship)
    user_id TEXT NOT NULL,

    -- Reference to the role assigned to the user
    role_id TEXT NOT NULL,

    -- Record creation timestamp, set automatically
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Composite primary key ensures uniqueness of (user_id, role_id) pairs
    PRIMARY KEY (user_id, role_id),

    -- When a user is deleted, automatically remove all their role assignments
    FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE,

    -- Reference to roles table; prevents deletion of a role if it is still assigned
    FOREIGN KEY (role_id)
        REFERENCES roles(role_id)
);

-- Index to speed up lookups by role_id
CREATE INDEX user_roles_role_id_idx
ON user_roles(role_id);