-- +goose Up
ALTER TABLE databases ADD COLUMN role VARCHAR(20) DEFAULT 'PRIMARY' NOT NULL;
ALTER TABLE databases ADD COLUMN primary_id UUID REFERENCES databases(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE databases DROP COLUMN primary_id;
ALTER TABLE databases DROP COLUMN role;
