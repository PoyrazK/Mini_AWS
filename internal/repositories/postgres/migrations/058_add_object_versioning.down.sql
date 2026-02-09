-- +goose Up
-- Handled in .up.sql

-- +goose Down
DROP INDEX IF EXISTS idx_objects_latest;
ALTER TABLE objects DROP CONSTRAINT IF EXISTS objects_bucket_key_version_unique;
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'objects_bucket_key_key') THEN
        ALTER TABLE objects ADD CONSTRAINT objects_bucket_key_key UNIQUE (bucket, key);
    END IF;
END $$;
ALTER TABLE objects DROP COLUMN IF EXISTS version_id;
ALTER TABLE objects DROP COLUMN IF EXISTS is_latest;
