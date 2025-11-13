-- Change news.image from BYTEA to TEXT path and keep existing values as UTF-8 strings
-- This assumes the existing BYTEA contains a URL or path encoded in UTF-8.
-- If the column is already TEXT/VARCHAR, this is a no-op in dev environments.

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='news' AND column_name='image' AND data_type='bytea'
    ) THEN
        ALTER TABLE news
        ALTER COLUMN image TYPE TEXT USING convert_from(image, 'UTF8');
    END IF;
END$$;

-- Optional: widen to VARCHAR if you prefer length constraint
-- ALTER TABLE news ALTER COLUMN image TYPE VARCHAR(512);
