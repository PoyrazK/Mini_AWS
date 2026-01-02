-- Migration: 012_add_failure_tracking.down.sql

ALTER TABLE scaling_groups DROP COLUMN IF EXISTS failure_count;
ALTER TABLE scaling_groups DROP COLUMN IF EXISTS last_failure_at;
