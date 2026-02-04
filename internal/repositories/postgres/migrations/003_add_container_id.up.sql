-- +goose Up

ALTER TABLE instances ADD COLUMN IF NOT EXISTS container_id TEXT;
