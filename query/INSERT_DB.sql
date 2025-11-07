-- Example user insert with a pre-hashed password (replace hash accordingly)
-- Password: admin123
-- Generate with: bcrypt (cost 10). Example value below is a placeholder; replace in production.
INSERT INTO users (username, email, password_hash, role, is_active, is_logged_in)
VALUES ('admin', 'admin@example.com', '$2a$10$cyWm6rTt6npl8mJ5lB6sEOe6iZfJbJr7Lz1mFQz2h6R0o1m0b8mwy', 'admin', true, false)
ON CONFLICT (username) DO NOTHING;

-- Seed categories
INSERT INTO categories (name, description) VALUES
 ('national','National news and policy'),
 ('international','Global news and events'),
 ('sports','Sports updates and analysis'),
 ('entertainment','Movies, music and culture'),
 ('technology','Tech and science innovations')
ON CONFLICT (name) DO NOTHING;

-- Seed example published news (store image as URL in BYTEA using decode)
-- Using convert_to to store URL text as bytea
WITH c AS (
	SELECT id, name FROM categories
)
INSERT INTO news (category_id, author_id, title, content, image, status, published_at)
SELECT (SELECT id FROM c WHERE name='national'), 1,
	'Government Announces New Economic Plan to Boost Job Growth',
	'The government unveiled a comprehensive plan aimed at boosting job growth across key sectors...',
	convert_to('https://cdn.usegalileo.ai/sdxl10/c211285b-0fdc-46ee-9467-9471799200dd.png','UTF8'),
	'published'::news_status, NOW()
UNION ALL
SELECT (SELECT id FROM c WHERE name='technology'), 1,
	'Tech Giant Unveils Latest Smartphone Model with Advanced Features',
	'A major tech company introduced its latest flagship smartphone featuring advanced AI capabilities...',
	convert_to('https://cdn.usegalileo.ai/sdxl10/5adb2e21-9930-44df-b4ec-db97de1aabce.png','UTF8'),
	'published'::news_status, NOW()
ON CONFLICT DO NOTHING;
