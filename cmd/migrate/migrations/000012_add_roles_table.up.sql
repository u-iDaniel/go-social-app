CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level int NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO roles (name, level, description) VALUES
('admin', 3, 'Administrator with full access can update and delete other posts'),
('moderator', 2, 'Moderator with limited access can update other posts'),
('user', 1, 'Regular user with basic access their own posts and comments');