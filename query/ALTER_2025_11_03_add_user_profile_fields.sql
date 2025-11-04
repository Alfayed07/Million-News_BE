-- Migration: Add profile fields to users table (non-destructive)
-- Date: 2025-11-03
-- Safe to run multiple times due to IF NOT EXISTS

BEGIN;

ALTER TABLE IF EXISTS users
    ADD COLUMN IF NOT EXISTS name   VARCHAR(100);

ALTER TABLE IF EXISTS users
    ADD COLUMN IF NOT EXISTS bio    TEXT;

ALTER TABLE IF EXISTS users
    ADD COLUMN IF NOT EXISTS avatar VARCHAR(255);

-- Optional backfill to empty strings (uncomment if you prefer no NULLs)
-- UPDATE users SET name = ''   WHERE name   IS NULL;
-- UPDATE users SET bio = ''    WHERE bio    IS NULL;
-- UPDATE users SET avatar = '' WHERE avatar IS NULL;

COMMIT;
