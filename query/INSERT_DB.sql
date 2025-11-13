-- Example user insert with a pre-hashed password (replace hash accordingly)
-- Password: admin123
-- Generate with: bcrypt (cost 10). Example value below is a placeholder; replace in production.
INSERT INTO users (username, email, password_hash, role, is_active, is_logged_in)
VALUES ('admin', 'admin@example.com', '$2a$10$cyWm6rTt6npl8mJ5lB6sEOe6iZfJbJr7Lz1mFQz2h6R0o1m0b8mwy', 'admin', true, false)
ON CONFLICT (username) DO NOTHING;

-- Second seeded editor user
-- Password: editor123 (bcrypt cost 10) -> example hash below
INSERT INTO users (username, email, password_hash, role, is_active, is_logged_in)
VALUES ('editor', 'editor@example.com', '$2a$10$5XkXpLxZCwV6xBIYgZtPJe1jYwQzIhJrXK2JxqgV7zBqQXoJ5nQ3e', 'editor', true, false)
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
-- Seed 10 news items using relative image paths (ensure files exist under public/uploads/)
WITH c AS (SELECT id, name FROM categories)
INSERT INTO news (category_id, author_id, title, content, image, status, published_at)
VALUES
((SELECT id FROM c WHERE name='national'), 1,
 'Pemerintah Resmikan Proyek Infrastruktur Nasional Skala Besar',
 'Pemerintah meresmikan proyek infrastruktur terintegrasi yang akan meningkatkan konektivitas antar daerah selama dekade mendatang...',
 '/uploads/national-1.jpg','published', NOW() - INTERVAL '2 hours'),
((SELECT id FROM c WHERE name='national'), 1,
 'Program Subsidi Pertanian Diperluas untuk Mendukung Petani Kecil',
 'Subsidi baru difokuskan pada peningkatan produktivitas pertanian berkelanjutan dan pemberdayaan petani kecil di berbagai provinsi...',
 '/uploads/national-2.jpg','published', NOW() - INTERVAL '5 hours'),
((SELECT id FROM c WHERE name='international'), 1,
 'Pertemuan Tingkat Tinggi Bahas Kerjasama Iklim Global',
 'Para pemimpin dunia berkumpul untuk menyepakati langkah konkret dalam pengurangan emisi karbon dan transisi energi bersih...',
 '/uploads/international-1.jpg','published', NOW() - INTERVAL '1 day'),
((SELECT id FROM c WHERE name='international'), 1,
 'Bantuan Kemanusiaan Dikirim ke Wilayah Terdampak Bencana Alam',
 'Organisasi internasional mengoordinasikan pengiriman logistik dan tenaga medis ke wilayah yang baru saja dilanda bencana...',
 '/uploads/international-2.jpg','published', NOW() - INTERVAL '3 days'),
((SELECT id FROM c WHERE name='sports'), 1,
 'Tim Nasional Mencetak Kemenangan Dramatis di Babak Akhir',
 'Gol di menit-menit terakhir memastikan kemenangan penting yang membuka peluang lolos ke turnamen regional...',
 '/uploads/sports-1.jpg','published', NOW() - INTERVAL '6 hours'),
((SELECT id FROM c WHERE name='sports'), 1,
 'Kompetisi Atlet Muda Menarik Perhatian Banyak Sponsor',
 'Ajang pencarian bakat olahraga tingkat nasional menarik dukungan luas dari sektor swasta...',
 '/uploads/sports-2.jpg','published', NOW() - INTERVAL '14 hours'),
((SELECT id FROM c WHERE name='entertainment'), 1,
 'Film Lokal Raih Penghargaan di Festival Internasional',
 'Sebuah film produksi dalam negeri memenangkan kategori sinematografi dan skenario terbaik di ajang festival bergengsi...',
 '/uploads/entertainment-1.jpg','published', NOW() - INTERVAL '20 hours'),
((SELECT id FROM c WHERE name='entertainment'), 1,
 'Musisi Muda Meluncurkan Album Debut Bernuansa Eksperimental',
 'Album perdana ini memadukan elemen tradisional dan elektronik yang mendapat respon positif dari kritikus musik...',
 '/uploads/entertainment-2.jpg','published', NOW() - INTERVAL '2 days'),
((SELECT id FROM c WHERE name='technology'), 1,
 'Startup AI Luncurkan Platform Analitik Berbasis Pembelajaran Mesin',
 'Platform baru membantu perusahaan memproses data real-time untuk keputusan yang lebih cepat dan akurat...',
 '/uploads/technology-1.jpg','published', NOW() - INTERVAL '4 hours'),
((SELECT id FROM c WHERE name='technology'), 1,
 'Peneliti Mengembangkan Material Baterai Ramah Lingkungan Generasi Baru',
 'Material ini diklaim memiliki densitas energi lebih tinggi dan proses produksi yang rendah emisi...',
 '/uploads/technology-2.jpg','published', NOW() - INTERVAL '3 days')
ON CONFLICT DO NOTHING;

-- Placeholder note: ensure sample image files exist (or use placeholder) under public/uploads.
-- You can upload new images via the /manage/upload endpoint once implemented.
