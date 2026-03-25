INSERT INTO roles (role_id, name, description) VALUES
('admin',       'Administrator', 'Full access to all system features and settings'),
('user',        'User',          'Regular user with standard access'),
('guest',       'Guest',         'Limited access, read-only or temporary sessions'),
('viewer_all',  'Viewer All',    'Can view all from all users');

-- Assign "guest" role to all users with is_guest = 1
INSERT INTO user_roles (user_id, role_id)
SELECT user_id, 'guest'
FROM users
WHERE is_guest = 1
  AND user_id NOT IN (
      SELECT user_id
      FROM user_roles
      WHERE role_id = 'guest'
  );