-- Rollback migration
ALTER TABLE gateway_routes 
    DROP COLUMN IF EXISTS pattern_type,
    DROP COLUMN IF EXISTS path_pattern,
    DROP COLUMN IF EXISTS param_names,
    DROP COLUMN IF EXISTS priority;

DROP INDEX IF EXISTS idx_gateway_routes_pattern_type;
