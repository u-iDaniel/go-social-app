ALTER TABLE IF EXISTS users
ADD COLUMN IF NOT EXISTS role_id INT REFERENCES roles(id);

UPDATE users SET role_id = (SELECT id FROM roles WHERE name = 'user') WHERE role_id IS NULL; -- Set default role for existing users to be 'user'

ALTER TABLE users ALTER COLUMN role_id SET NOT NULL;