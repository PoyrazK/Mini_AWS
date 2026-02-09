-- +goose Up
ALTER TABLE scaling_groups ADD COLUMN IF NOT EXISTS failure_count INT DEFAULT 0;
ALTER TABLE scaling_groups ADD COLUMN IF NOT EXISTS last_failure_at TIMESTAMPTZ;
