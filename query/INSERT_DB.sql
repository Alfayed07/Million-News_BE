-- Example user insert with a pre-hashed password (replace hash accordingly)
-- Password: admin123
-- Generate with: bcrypt (cost 10). Example value below is a placeholder; replace in production.
INSERT INTO users (username, email, password_hash, role, is_active, is_logged_in)
VALUES ('admin', 'admin@example.com', '$2a$10$cyWm6rTt6npl8mJ5lB6sEOe6iZfJbJr7Lz1mFQz2h6R0o1m0b8mwy', 'admin', true, false)
ON CONFLICT (username) DO NOTHING;
