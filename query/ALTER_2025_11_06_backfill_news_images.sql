-- Backfill news.image (BYTEA) for missing images
-- Assumes the column `news.image` is BYTEA storing URL text as UTF-8 bytes.
-- This script sets stable placeholder URLs only when image is NULL or empty.

BEGIN;

UPDATE public.news SET image = convert_to('https://picsum.photos/seed/news-3/1200/630','UTF8')
WHERE id = 3 AND (image IS NULL OR octet_length(image) = 0);

UPDATE public.news SET image = convert_to('https://picsum.photos/seed/news-4/1200/630','UTF8')
WHERE id = 4 AND (image IS NULL OR octet_length(image) = 0);

UPDATE public.news SET image = convert_to('https://picsum.photos/seed/news-5/1200/630','UTF8')
WHERE id = 5 AND (image IS NULL OR octet_length(image) = 0);

UPDATE public.news SET image = convert_to('https://picsum.photos/seed/news-6/1200/630','UTF8')
WHERE id = 6 AND (image IS NULL OR octet_length(image) = 0);

UPDATE public.news SET image = convert_to('https://picsum.photos/seed/news-7/1200/630','UTF8')
WHERE id = 7 AND (image IS NULL OR octet_length(image) = 0);

UPDATE public.news SET image = convert_to('https://picsum.photos/seed/news-8/1200/630','UTF8')
WHERE id = 8 AND (image IS NULL OR octet_length(image) = 0);

UPDATE public.news SET image = convert_to('https://picsum.photos/seed/news-9/1200/630','UTF8')
WHERE id = 9 AND (image IS NULL OR octet_length(image) = 0);

UPDATE public.news SET image = convert_to('https://picsum.photos/seed/news-10/1200/630','UTF8')
WHERE id = 10 AND (image IS NULL OR octet_length(image) = 0);

COMMIT;
