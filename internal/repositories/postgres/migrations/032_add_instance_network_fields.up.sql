-- +goose Up
ALTER TABLE instances ADD COLUMN IF NOT EXISTS subnet_id UUID REFERENCES subnets(id);
ALTER TABLE instances ADD COLUMN IF NOT EXISTS private_ip INET;
ALTER TABLE instances ADD COLUMN IF NOT EXISTS ovs_port VARCHAR(64);

CREATE INDEX IF NOT EXISTS idx_instances_subnet ON instances(subnet_id);
